#/usr/bin/env bash

# Usage ./get.sh <token>

curl --header "Content-Type: application/json" \
  --request GET \
  --data '{"token":"'"$1"'"}' \
  http://localhost:8080/upload
