#!/usr/bin/env bash

# Usage: ./register.sh <login> <password>

curl -X POST localhost:8080/register \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "login=$1&password=$2"
