#!/usr/bin/env bash

# Usage: ./upload.sh <token> <file>

curl -X GET \
  -F "token=$1" \
  http://localhost:8080/upload


