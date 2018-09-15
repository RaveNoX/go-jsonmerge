#!/bin/sh

set -e

go get -u github.com/golang/dep/cmd/dep
dep ensure -v
