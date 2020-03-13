package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/RaveNoX/go-jsonmerge"
)

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

	fmt.Printf("%v:\n", file)
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

func printPatch(file string, info *jsonmerge.Info) {
	m := make(map[string]interface{})

	for k, v := range info.Replaced {
		parts := strings.Split(k, ".")
		lastIndex := len(parts) - 1
		keyObject := m
		for _, p := range parts[:lastIndex] {
			if o, ok := keyObject[p].(map[string]interface{}); ok {
				keyObject = o
			} else {
				o = make(map[string]interface{})
				keyObject[p] = o
				keyObject = o
			}
		}

		keyObject[parts[lastIndex]] = v
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	fmt.Printf("%v:\n", file)
	err := enc.Encode(m)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot write %s patch: %v\n", file, err)
	}

	fmt.Println()
}
