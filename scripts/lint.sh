#!/usr/bin/env bash

set -o errexit
set -o nounset

if ! command -v golangci-lint &> /dev/null
then
    echo "golangci-lint could not be found"
    exit
fi

golangci-lint run ./...

echo "no linting problems found"
