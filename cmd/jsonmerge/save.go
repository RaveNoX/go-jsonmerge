package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func save(filePath string, data interface{}) error {
	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("Cannot create file: %v", err)
	}
	defer file.Close()

	err = saveJSON(file, data)
	if err != nil {
		return fmt.Errorf("Cannot write JSON: %v", err)
	}

	return nil
}

func saveJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
