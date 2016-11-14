package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/RaveNoX/go-jsoncommentstrip"
	"github.com/spkg/bom"
)

func load(filePath string) (interface{}, error) {
	if filePath == "-" {
		fmt.Fprintln(os.Stderr, "Loading patch from STDIN")
		return loadFromReader(os.Stdin)
	}

	return loadFromFile(filePath)
}

func loadFromReader(reader io.Reader) (interface{}, error) {
	var data interface{}

	bomReader := bom.NewReader(reader)
	jsonCommentReader := jsoncommentstrip.NewReader(bomReader)

	dec := json.NewDecoder(jsonCommentReader)
	dec.UseNumber()

	err := dec.Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("Cannot parse JSON: %v", err)
	}

	return data, nil
}

func loadFromFile(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("Cannot open file: %v", err)
	}
	defer file.Close()

	return loadFromReader(file)
}
