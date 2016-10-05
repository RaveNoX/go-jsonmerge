package jsonmerge

import (
	"encoding/json"
	"testing"
)

var (
	outputJSON, outputIndentJSON string
	input                        = `
{  
  "number": 1,
  "string": "value",
  "object": {
    "number": 1,
    "string": "value",
    "nested_object": {
      "number": 2
    },
    "array": [1, 2, 3],
    "partial_array": [1, 2, 3]
  }
}
    `
	patch = `
{  
  "number": 2,
  "string": "value1",
  "nonexitent": "woot",
  "object": {
    "number": 3,
    "string": "value2",
    "nested_object": {
      "number": 4
    },
    "array": [3, 2, 1],
    "partial_array": {
      "1": 4
    }
  }
}
    `
)

func init() {
	output := []byte(`
{  
  "number": 2,
  "string": "value1",
  "object": {
    "number": 3,
    "string": "value2",
    "nested_object": {
      "number": 4
    },
    "array": [3, 2, 1],
    "partial_array": [1, 4, 3]
  }
}
    `)

	var outputData interface{}
	json.Unmarshal(output, &outputData)

	output, _ = json.Marshal(outputData)
	outputJSON = string(output)

	output, _ = json.MarshalIndent(outputData, " ", "  ")
	outputIndentJSON = string(output)
}

func TestMergeBytesIndent(t *testing.T) {
	result, info, err := MergeBytesIndent([]byte(input), []byte(patch), " ", "  ")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(result) != outputIndentJSON {
		t.Errorf("Result not equals output\nExpected:\n%s\n\nGot:\n%s\n\n", outputIndentJSON, result)
	}

	if len(info.Errors) != 0 {
		t.Errorf("info.Errors count is not 0, count: %v", len(info.Errors))
	}
}

func TestMergeBytes(t *testing.T) {
	result, info, err := MergeBytes([]byte(input), []byte(patch))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(result) != outputJSON {
		t.Errorf("Result not equals output\nExpected:\n%s\n\nGot:\n%s\n\n", outputJSON, result)
	}

	if len(info.Errors) != 0 {
		t.Errorf("info.Errors count is not 0, count: %v", len(info.Errors))
	}
}

func TestLongNumbers(t *testing.T) {
	input := `{"Id":12423434,"Value":12423434}`
	patch := `{"Value":12423439}`
	outputJSON := `{"Id":12423434,"Value":12423439}`

	result, info, err := MergeBytes([]byte(input), []byte(patch))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(result) != outputJSON {
		t.Errorf("Result not equals output\nExpected:\n%s\n\nGot:\n%s\n\n", outputJSON, result)
	}

	if len(info.Errors) != 0 {
		t.Errorf("info.Errors count is not 0, count: %v", len(info.Errors))
	}
}
