#!/bin/env bash
set -e

if [ -z "$GH_BOT_TOKEN" ]; then
    echo "GH_BOT_TOKEN is missing, abort."
    exit 1
fi

# https://docs.github.com/en/free-pro-team@latest/rest/reference/pulls#create-a-pull-request
result=$(curl -s -X POST -H "Authorization: token $GH_BOT_TOKEN" -d "{\"head\":\"autobuild-Documentation-$TAG\",\"base\":\"master\",\"title\":\"Documentation $TAG\",\"body\":\"Automatic generation of the documentation $TAG\"}" "https://api.github.com/repos/outscale/terraform-provider-outscale/pulls")

errors=$(echo $result | jq .errors)

if [ "$errors" != "null" ]; then
    echo "errors while creating pull request, abort."
    exit 1
fi
