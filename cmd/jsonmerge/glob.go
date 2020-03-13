package main

import (
	"fmt"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
)

func getGlobFiles(globs []string) (files []string, err error) {
	var globFiles = make(map[string]interface{})
	var matches []string

	for _, glob := range globs {
		glob = filepath.FromSlash(glob)
		matches, err = doublestar.Glob(glob)

		if err != nil {
			err = fmt.Errorf("Glob error \"%v\": %v", glob, err)
			return
		}

		for _, match := range matches {
			globFiles[match] = true
		}
	}

	for file := range globFiles {
		files = append(files, file)
	}

	return
}
