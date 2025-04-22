#!/usr/bin/env bash

# Usage: ./login.sh <login> <password>

curl \
  -X POST \
  -H "Accept: application/json" \
  -d "login=$1&password=$2" \
  localhost:8080/login 
