#!/usr/bin/env bash

# Usage ./logout.sh <token>

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"token":"'"$1"'"}' \
  http://localhost:8080/logout
