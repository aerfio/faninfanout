package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

var fileNumber = 2000

func main() {
	log := logrus.New()

	bodies := make(chan string, fileNumber)
	_, executableFile, _, _ := runtime.Caller(0)
	//fmt.Println(executableFile)
	eg := errgroup.Group{}

	for i := 1; i <= fileNumber; i++ {
		fileName := filepath.Join(filepath.Dir(executableFile), "yamls", fmt.Sprintf("deploy-%d.yaml", i))
		eg.Go(func() error {
			log.Infof("start processing %s", fileName)
			json, err := process(fileName)
			if err != nil {
				return err
			}
			bodies <- json
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatalf("while waiting for goroutines to finish: %s", err)
	}

	close(bodies)
	for msg := range bodies {
		log.Info(msg)
	}
}

func process(fileName string) (string, error) {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("while reading file %s: %s", fileName, err)
	}

	jsonBytes, err := yamlutil.ToJSON(fileBytes)
	if err != nil {
		return "", fmt.Errorf("while casting yaml to json for file %s: %s", fileName, err)
	}

	time.Sleep(3 * time.Second) // oh no, very expensive task
	return gjson.GetBytes(jsonBytes, "metadata.name").String(), nil
}
