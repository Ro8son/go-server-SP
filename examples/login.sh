#!/usr/bin/env bash

# Usage ./login.sh <login> <password>

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"login":"'"$1"'","password":"'"$2"'"}' \
  http://ec2-13-60-9-150.eu-north-1.compute.amazonaws.com:8000/login
  #http://localhost:8000/login
