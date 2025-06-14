#!/usr/bin/env bash

TMP_JSON=$(mktemp)
TOKEN="$1"
FILE="$2"
URL="$3"
COUNT="$4"


cat > "$TMP_JSON" <<EOF
{
  "token": "$TOKEN",
  "file_id": $FILE,
  "url": "$URL",
  "max_uses": $COUNT
}
EOF

curl -X POST "localhost:8080/file/share/add" \
  -H "Content-Type: application/json" \
  -d @"$TMP_JSON"

rm "$TMP_JSON"
