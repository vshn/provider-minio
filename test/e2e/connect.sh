#!/bin/bash

set -eo pipefail
#set -x

endpoint=${1}
secret_name=${2}

access_key=$(kubectl -n default get secret ${secret_name} -o jsonpath='{.data.AWS_ACCESS_KEY_ID}' | base64 -d)
secret_key=$(kubectl -n default get secret ${secret_name} -o jsonpath='{.data.AWS_SECRET_ACCESS_KEY}' | base64 -d)
export MC_HOST_minio=http://${access_key}:${secret_key}@${endpoint}

${GOBIN}/mc mb --quiet "minio/bucket-$(date +%s)"
