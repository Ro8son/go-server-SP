#!/usr/bin/env bash

# Usage ./get_file_list.sh <token>

curl -X GET \
  -F "token=$1" \
  http://localhost:8080/file/list
