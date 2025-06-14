#!/usr/bin/env bash

TMP_JSON=$(mktemp)

cat > "$TMP_JSON" <<EOF
{
  "token": "'"$1"'"
}
EOF

curl -X POST "localhost:8080/album/add" \
  -H "Content-Type: application/json" \
  -d @"$TMP_JSON"

rm "$TMP_JSON"
