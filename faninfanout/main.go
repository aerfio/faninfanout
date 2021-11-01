package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

var fileNumber = 2000
var workerNum = 500

// https://gobyexample.com/worker-pools
func main() {
	log := logrus.New()

	bodies := make(chan string, fileNumber)
	errs := make(chan error, fileNumber)
	jobs := make(chan string, fileNumber)
	_, executableFile, _, _ := runtime.Caller(0)
	//fmt.Println(filename)



	// start workers
	for i := 1; i <= workerNum; i++ {
		log.Infof("starting worker %d", i)
		go worker(jobs, bodies, errs)
	}

	// send work
	for i := 1; i <= fileNumber; i++ {
		fileName := filepath.Join(filepath.Dir(executableFile), "..", "yamls", fmt.Sprintf("deploy-%d.yaml", i))
		if i%250 == 0 {
			log.Infof("sending work %s", fileName)
		}
		jobs <- fileName
	}
	close(jobs)

	for i := 1; i <= fileNumber; i++ {
		err := <-errs
		if err != nil {
			log.Fatalf("while receiving errors from error channel: %s", err)
		}
	}
	close(errs) // close errors after we've ensured that no more worker is sending there

	close(bodies)
	for msg := range bodies {
		log.Info(msg)
	}
}

func process(fileName string) (string, error) {
	logrus.Infof("start processing %s", fileName)
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("while reading file %s: %s", fileName, err)
	}

	jsonBytes, err := yamlutil.ToJSON(fileBytes)
	if err != nil {
		return "", fmt.Errorf("while casting yaml to json for file %s: %s", fileName, err)
	}

	time.Sleep(3 * time.Second) // oh no, very expensive task
	logrus.Infof("finish processing %s", fileName)
	return gjson.GetBytes(jsonBytes, "metadata.name").String(), nil
}

func worker(jobs <-chan string, results chan<- string, errs chan<- error) {
	for j := range jobs {
		out, err := process(j)
		errs <- err
		results <- out
	}
}
