#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail
set -o errtrace

ROOT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." &> /dev/null && pwd )"

for i in {1..2000}
do
   # echo "kubectl create deployment nginx-${i} --image nginx --dry-run=client -oyaml > ${ROOT_DIR}/yamls/deploy-${i}.yaml" # comment it, it will be faster
   kubectl create deployment "nginx-${i}" --image nginx --dry-run=client -oyaml > "${ROOT_DIR}/yamls/deploy-${i}.yaml"
done
