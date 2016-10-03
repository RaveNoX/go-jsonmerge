package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/RaveNoX/go-jsoncommentstrip"
	"github.com/RaveNoX/go-jsonmerge"
	"github.com/bmatcuk/doublestar"
	"github.com/spkg/bom"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v <patch.json> <glob1.json> <glob2.json>...<globN.json>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	patchFile := os.Args[1]
	globs := os.Args[2:]

	patchBuff, err := readJSON(patchFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read patch file: %v\n", err)
		os.Exit(2)
	}

	dataFiles, err := getGlobFiles(globs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get data files: %v\n", err)
		os.Exit(3)
	}

	patchFiles(patchBuff, dataFiles, true)
}

func patchFiles(patchBuff []byte, dataFiles []string, replaces bool) {
	dataFilesCount := len(dataFiles)

	fmt.Printf("Data files to patch: %v\n\n", dataFilesCount)

	for _, file := range dataFiles {
		fmt.Printf("%v:\n", file)

		buff, err := readJSON(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot load data file \"%v\": %v\n", file, err)
			continue
		}

		result, info, err := jsonmerge.MergeBytesIndent(buff, patchBuff, "", "  ")

		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot merge data file \"%v\": %v\n", file, err)
			continue
		}

		if len(info.Errors) > 0 {
			fmt.Fprintf(os.Stderr, "Replacement errors for data file \"%v\":\n", file)
			for _, err := range info.Errors {
				fmt.Fprintf(os.Stderr, "  %v\n", err)
			}
			fmt.Fprintln(os.Stderr)
		}

		fmt.Printf("  Replaced %v items\n", len(info.Replaced))
		if replaces {
			for k, v := range info.Replaced {
				vBuff, _ := json.Marshal(v)

				fmt.Printf("    %v => %s\n", k, vBuff)
			}
			fmt.Println()
		}

		err = ioutil.WriteFile(file, result, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write data file \"%v\": %v\n", file, err)
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
