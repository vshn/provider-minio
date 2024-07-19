#!/bin/bash

set -eo pipefail

endpoint=${1}
bucket_name=${2}
file_path=${3}
secret_name=${4}
secret_namespace=${5}

access_key=$(kubectl -n "${secret_namespace}" get secret "${secret_name}" -o jsonpath='{.data.AWS_ACCESS_KEY_ID}' | base64 -d)
secret_key=$(kubectl -n "${secret_namespace}" get secret "${secret_name}" -o jsonpath='{.data.AWS_SECRET_ACCESS_KEY}' | base64 -d)
export MC_HOST_minio=http://${access_key}:${secret_key}@${endpoint}

echo "Uploading object to bucket: ${bucket_name}"
echo "File path: ${file_path}"
echo "Endpoint: ${endpoint}"
echo "Access key: ${access_key}"
echo "Secret key: ${secret_key}"

"${GOBIN}/mc" cp --quiet --debug "${file_path}" "minio/${bucket_name}"
