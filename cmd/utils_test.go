package main

import "io/ioutil"

func createTempDir() (string, error) {
	return ioutil.TempDir("test", "jsonmerge_test")
}
