#!/usr/bin/env bash

# Usage ./login.sh <login> <password>

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"token":"'"$1"'"}' \
  http://localhost:8080/file/list
