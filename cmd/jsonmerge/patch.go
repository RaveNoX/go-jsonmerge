package main

import (
	"fmt"
	"os"

	"github.com/RaveNoX/go-jsonmerge"
)

func patchFile(patchBuff interface{}, file string) error {
	data, err := loadFromFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: load error: %v\n", file, err)
		return fmt.Errorf("%v: load error: %v", file, err)
	}

	result, info := jsonmerge.Merge(data, patchBuff)

	if len(info.Errors) > 0 {
		printMergeErrors(file, info)
		return fmt.Errorf("%v: got %v merge errors", file, len(info.Errors))
	}

	if !opts.DryRun {
		err = save(file, result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: save error: %v\n", file, err)
			return fmt.Errorf("%v: save error: %v", file, err)
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

	return nil
}

func patchFiles(patchBuff interface{}, dataFiles []string) (errors []error) {
	for _, file := range dataFiles {
		err := patchFile(patchBuff, file)

		if err != nil {
			errors = append(errors, err)
		}
	}

	return
}
