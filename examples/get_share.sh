#!/usr/bin/env bash

TMP_JSON=$(mktemp)
TOKEN="$1"

cat > "$TMP_JSON" <<EOF
{
  "token": "$TOKEN"
}
EOF

curl -X POST "localhost:8080/file/share/get" \
  -H "Content-Type: application/json" \
  -d @"$TMP_JSON"

rm "$TMP_JSON"
