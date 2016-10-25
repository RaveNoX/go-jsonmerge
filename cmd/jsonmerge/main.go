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

	"sort"

	"github.com/bmatcuk/doublestar"
	"github.com/spkg/bom"
)

var opts *options.Options

func init() {
	opts = &options.Options{
		Name: filepath.Base(os.Args[0]),
	}
}

func main() {
	opts.ParseOrExit(os.Args[1:])

	patchBuff, err := readJSON(opts.Patch)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Patch load error: %v\n", err)
		os.Exit(2)
	}

	dataFiles, err := getGlobFiles(opts.Globs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get data files: %v\n", err)
		os.Exit(3)
	}

	if opts.DryRun {
		fmt.Fprintln(os.Stderr, "Dry run: no files will be really patched")
		fmt.Fprintln(os.Stderr)
	}
	patchFiles(patchBuff, dataFiles)
}

func printMergeErrors(file string, info *jsonmerge.Info) {
	for _, err := range info.Errors {
		fmt.Fprintf(os.Stderr, "%v: replace warning: %v\n", file, err)
	}
	fmt.Fprintln(os.Stderr)
}

func printReplaces(file string, info *jsonmerge.Info) {
	var keys []string

	for k := range info.Replaced {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := info.Replaced[k]
		vBuff, err := json.Marshal(v)
		if err != nil {
			vBuff = []byte(fmt.Sprintf("<cannot get value as JSON: %v>", err))
		}

		fmt.Printf("  %v = %s\n", k, vBuff)
	}
	fmt.Println()
}

func patchFile(patchBuff []byte, file string) {
	buff, err := readJSON(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: load error: %v\n", file, err)
		return
	}

	result, info, err := jsonmerge.MergeBytesIndent(buff, patchBuff, "", "  ")

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: merge error: %v\n", file, err)
		return
	}

	if len(info.Errors) > 0 {
		printMergeErrors(file, info)
	}

	if !opts.DryRun {
		err = ioutil.WriteFile(file, result, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: save error: %v\n", file, err)
		}
	}

	if !opts.Quiet {
		if opts.Verbose {
			fmt.Printf("%v:\n", file)
			printReplaces(file, info)
		} else {
			fmt.Printf("%v\n", file)
		}
	}
}

func patchFiles(patchBuff []byte, dataFiles []string) {
	for _, file := range dataFiles {
		patchFile(patchBuff, file)
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
