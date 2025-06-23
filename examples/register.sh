#!/usr/bin/env bash

# Usage ./register.sh <login> <password>


curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"login":"'"$1"'","password":"'"$2"'"}' \
  http://localhost:8000/register
