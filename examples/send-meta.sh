#/usr/bin/env bash

# Usage ./send-meta.sh <token>

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"token":"'"$1"'","file":"'"$( base64 ~/pictures/cute/166a71ed924501029b1af13e48ad804b.png)"'"}' \
  http://localhost:8080/metadata
