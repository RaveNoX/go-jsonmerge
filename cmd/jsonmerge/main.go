package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RaveNoX/go-jsonmerge/cmd/jsonmerge/options"
)

var opts *options.Options

func init() {
	opts = &options.Options{
		Name: filepath.Base(os.Args[0]),
	}
}

func main() {
	opts.ParseOrExit(os.Args[1:])

	patchBuff, err := load(opts.PatchFile)

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

	if errors := patchFiles(patchBuff, dataFiles); len(errors) > 0 {
		fmt.Fprintln(os.Stderr, "Errors:")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "    %v\n", err)
		}

		os.Exit(4)
	}
}
