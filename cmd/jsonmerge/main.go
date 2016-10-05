package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/RaveNoX/go-jsonmerge"
	"github.com/RaveNoX/go-jsonmerge/cmd/jsonmerge/options"

	"github.com/RaveNoX/go-jsoncommentstrip"

	"github.com/bmatcuk/doublestar"
	"github.com/spkg/bom"
)

var settings *options.Options

func init() {
	settings = &options.Options{
		Name: filepath.Base(os.Args[0]),
	}
}

func main() {
	settings.ParseOrExit(os.Args[1:])

	patchBuff, err := readJSON(settings.Patch)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Patch load error: %v\n", err)
		os.Exit(2)
	}

	dataFiles, err := getGlobFiles(settings.Globs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get data files: %v\n", err)
		os.Exit(3)
	}

	patchFiles(patchBuff, dataFiles, true)
}

func patchFiles(patchBuff []byte, dataFiles []string, replaces bool) {
	for _, file := range dataFiles {
		buff, err := readJSON(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: load error: %v\n", file, err)
			continue
		}

		result, info, err := jsonmerge.MergeBytesIndent(buff, patchBuff, "", "  ")

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: merge error: %v\n", file, err)
			continue
		}

		if len(info.Errors) > 0 {
			for _, err := range info.Errors {
				fmt.Fprintf(os.Stderr, "%v: replace warning: %v\n", file, err)
			}
			fmt.Fprintln(os.Stderr)
		}

		err = ioutil.WriteFile(file, result, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: save error: %v\n", file, err)
		}

		if !settings.Quiet {
			if settings.Verbose {
				fmt.Printf("%v:\n", file)
				if replaces {
					for k, v := range info.Replaced {
						vBuff, _ := json.Marshal(v)

						fmt.Printf("  %v = %s\n", k, vBuff)
					}
					fmt.Println()
				}
			} else {
				fmt.Printf("%v\n", file)
			}
		}
	}
}

func getGlobFiles(globs []string) (files []string, err error) {
	var globFiles = make(map[string]interface{})
	var matches []string

	for _, glob := range globs {
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

func readJSON(filePath string) (buff []byte, err error) {
	file, err := os.Open(filePath)

	if err != nil {
		err = fmt.Errorf("Cannot open file: %v", err)
		return
	}
	defer file.Close()

	bomReader := bom.NewReader(file)
	jsonCommentReader := jsoncommentstrip.NewReader(bomReader)

	buff, err = ioutil.ReadAll(jsonCommentReader)

	if err != nil {
		err = fmt.Errorf("Cannot read JSON: %v", err)
	}
	return
}
