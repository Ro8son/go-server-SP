#/usr/bin/env bash

# Usage ./get.sh <token>

curl -X GET \
  -F "token=$1" \
  http://localhost:8080/file/list
