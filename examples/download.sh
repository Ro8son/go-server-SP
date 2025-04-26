#!/usr/bin/env bash

# Usage: ./upload.sh <token> <file>

curl -X GET \
  -F "token=$1" \
  -F "file_name=$2" \
  http://localhost:8080/file/download \
  --output $2


