#!/usr/bin/env bash

# Usage: ./upload.sh <token> <file>

curl -X POST \
  -F "file=@$2" \
  -F "token=$1" \
  http://localhost:8080/upload


