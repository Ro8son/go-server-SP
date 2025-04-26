#/usr/bin/env bash

# Usage ./send-meta.sh <token>

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"token":"'"$1"'"}' \
  http://localhost:8080/metadata
