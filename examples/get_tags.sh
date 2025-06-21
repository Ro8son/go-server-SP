#!/usr/bin/env bash

# Usage ./get_tags.sh <token>

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"token":"'"$1"'"}' \
  http://localhost:8000/file/tags
