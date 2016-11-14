package main

import "io/ioutil"
import "os"

func createTempDir() (string, error) {
	err := os.MkdirAll("test", os.ModeDir)

	if err != nil {
		return "", err
	}

	return ioutil.TempDir("test", "jsonmerge_test")
}
