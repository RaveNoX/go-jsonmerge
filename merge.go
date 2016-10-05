package jsonmerge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Info describes result of merge operation
type Info struct {
	// Errors is slice of non-critical errors of merge operations
	Errors []error
	// Replaced is describe replacements
	// Key is path in document like
	//   "prop1.prop2.prop3" for object properties or
	//   "arr1.1.prop" for arrays
	// Value is value of replacemet
	Replaced map[string]interface{}
}

func (info *Info) mergeValue(path []string, patch map[string]interface{}, key string, value interface{}) interface{} {
	patchValue, patchHasValue := patch[key]

	if !patchHasValue {
		return value
	}

	_, patchValueIsObject := patchValue.(map[string]interface{})

	path = append(path, key)
	pathStr := strings.Join(path, ".")

	if _, ok := value.(map[string]interface{}); ok {
		if !patchValueIsObject {
			err := fmt.Errorf("patch value must be object for key \"%v\"", pathStr)
			info.Errors = append(info.Errors, err)
			return value
		}

		return info.mergeObjects(value, patchValue, path)
	}

	if _, ok := value.([]interface{}); ok && patchValueIsObject {
		return info.mergeObjects(value, patchValue, path)
	}

	if !reflect.DeepEqual(value, patchValue) {
		info.Replaced[pathStr] = patchValue
	}

	return patchValue
}

func (info *Info) mergeObjects(data, patch interface{}, path []string) interface{} {
	if patchObject, ok := patch.(map[string]interface{}); ok {
		if dataArray, ok := data.([]interface{}); ok {
			ret := make([]interface{}, len(dataArray))

			for i, val := range dataArray {
				ret[i] = info.mergeValue(path, patchObject, strconv.Itoa(i), val)
			}

			return ret
		} else if dataObject, ok := data.(map[string]interface{}); ok {
			ret := make(map[string]interface{})

			for k, v := range dataObject {
				ret[k] = info.mergeValue(path, patchObject, k, v)
			}

			return ret
		}
	}

	return data
}

// Merge merges patch document to data document
//
// Returning merged document and merge info
func Merge(data, patch interface{}) (interface{}, *Info) {
	info := &Info{
		Replaced: make(map[string]interface{}),
	}
	ret := info.mergeObjects(data, patch, nil)
	return ret, info
}

// MergeBytesIndent merges patch document buffer to data document buffer
//
// Use prefix and indent for set indentation like in json.MarshalIndent
//
// Returning merged document buffer, merge info and
// error if any
func MergeBytesIndent(dataBuff, patchBuff []byte, prefix, indent string) (mergedBuff []byte, info *Info, err error) {
	var data, patch, merged interface{}

	err = unmarshalJSON(dataBuff, &data)
	if err != nil {
		err = fmt.Errorf("error in data JSON: %v", err)
		return
	}

	err = unmarshalJSON(patchBuff, &patch)
	if err != nil {
		err = fmt.Errorf("error in patch JSON: %v", err)
		return
	}

	merged, info = Merge(data, patch)

	mergedBuff, err = json.MarshalIndent(merged, prefix, indent)
	if err != nil {
		err = fmt.Errorf("error writing merged JSON: %v", err)
	}

	return
}

// MergeBytes merges patch document buffer to data document buffer
//
// Returning merged document buffer, merge info and
// error if any
func MergeBytes(dataBuff, patchBuff []byte) (mergedBuff []byte, info *Info, err error) {
	var data, patch, merged interface{}

	err = unmarshalJSON(dataBuff, &data)
	if err != nil {
		err = fmt.Errorf("error in data JSON: %v", err)
		return
	}

	err = unmarshalJSON(patchBuff, &patch)
	if err != nil {
		err = fmt.Errorf("error in patch JSON: %v", err)
		return
	}

	merged, info = Merge(data, patch)

	mergedBuff, err = json.Marshal(merged)
	if err != nil {
		err = fmt.Errorf("error writing merged JSON: %v", err)
	}

	return
}

func unmarshalJSON(buff []byte, data interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(buff))
	decoder.UseNumber()

	return decoder.Decode(data)
}
