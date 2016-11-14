package main

import (
	"io/ioutil"
	"os"
)

func createTempDir() (string, error) {
	return ioutil.TempDir(os.TempDir(), "jsonmerge_test")
}
