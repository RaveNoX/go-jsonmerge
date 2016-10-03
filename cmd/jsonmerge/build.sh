#!/bin/sh

set -e

GOARCH=adm64 GOOS=linux go build
GOARCH=adm64 GOOS=windows go build