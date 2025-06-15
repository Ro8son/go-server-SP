#!/usr/bin/env bash

# Usage: ./upload.sh <TOKEN> <FILE_NAME>
TOKEN="$1"
FILE_NAME="$2"
API_ENDPOINT="localhost:8000/file/upload"

if [[ -z "$TOKEN" || -z "$FILE_NAME" ]]; then
  echo "Usage: $0 <TOKEN> <FILE_NAME>"
  exit 1
fi

if [[ ! -f "$FILE_NAME" ]]; then
  echo "File '$FILE_NAME' does not exist."
  exit 1
fi

BASE64_CONTENT=$(base64 -w 0 "$FILE_NAME")
TMP_JSON=$(mktemp)

cat > "$TMP_JSON" <<EOF
{
  "token": "$TOKEN",
  "files": [
    {
      "file": "$BASE64_CONTENT",
      "metadata": {
        "file_name": "$FILE_NAME"
      }
    }
  ]
}
EOF

curl -X POST "$API_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d @"$TMP_JSON"

rm "$TMP_JSON"

