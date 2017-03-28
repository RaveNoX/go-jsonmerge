#!/bin/sh

set -e

MY_DIR=$(dirname $(readlink -f "$0"))

cd "${MY_DIR}"
mkdir -p "artifacts"

echo "Linux"
GOARCH=adm64 GOOS=linux go build -o "artifacts/jsonmerge" ./cmd

echo "Windows"
GOARCH=adm64 GOOS=windows go build -o "artifacts/jsonmerge.exe" ./cmd

echo "Build done"
