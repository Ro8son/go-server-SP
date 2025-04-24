#!/usr/bin/env bash

# Usage: ./upload.sh <token> <file>

curl -X POST \
  -F "token=$1" \
  -F "file=@$2" \
  -F "transaction_id=$3" \
  http://localhost:8080/upload


