package main

import (
	"fmt"
	"os"

	"github.com/RaveNoX/go-jsonmerge"
)

func patchFile(patchBuff interface{}, file string) {
	data, err := loadFromFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: load error: %v\n", file, err)
		return
	}

	result, info := jsonmerge.Merge(data, patchBuff)

	if len(info.Errors) > 0 {
		printMergeErrors(file, info)
	}

	if !opts.DryRun {
		err = save(file, result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: save error: %v\n", file, err)
		}
	}

	if !opts.Quiet {
		if opts.Verbose {
			printReplaces(file, info)
		} else if opts.Patch {
			printPatch(file, info)
		} else {
			fmt.Printf("%v\n", file)
		}
	}
}

func patchFiles(patchBuff interface{}, dataFiles []string) {
	for _, file := range dataFiles {
		patchFile(patchBuff, file)
	}
}
