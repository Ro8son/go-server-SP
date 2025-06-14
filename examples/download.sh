#!/usr/bin/env bash

# Usage: ./download.sh <TOKEN> <FILE_ID_1> [<FILE_ID_2> ...]
TOKEN="$1"
API_ENDPOINT="localhost:8080/file/download"

if [[ -z "$TOKEN" || $# -lt 2 ]]; then
  echo "Usage: $0 <TOKEN> <FILE_ID_1> [<FILE_ID_2> ...]"
  exit 1
fi

shift  # Remove token from arguments

# Convert arguments to integers and build a JSON array
FILE_IDS_JSON=$(printf '%s\n' "$@" | jq -R 'tonumber' | jq -s .)

TMP_JSON=$(mktemp)
cat > "$TMP_JSON" <<EOF
{
  "token": "$TOKEN",
  "file_ids": $FILE_IDS_JSON
}
EOF

curl -X POST "$API_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d @"$TMP_JSON"

rm "$TMP_JSON"


