package main

import "testing"
import "os"
import "io/ioutil"
import "path/filepath"
import "sort"

import "reflect"

func createGlobFiles(files []string) error {
	data := []byte(`
    {
        "test": 1
    }
    `)

	for _, item := range files {
		err := os.MkdirAll(filepath.Dir(item), os.ModeDir)

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(item, data, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestGlob(t *testing.T) {
	dir, err := createTempDir()

	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	files := []string{
		"test.json",
		"woot/test.json",
		"woot/abc/test.json",
		"gg/t/a.json",
	}

	for i := range files {
		files[i] = filepath.Join(dir, filepath.FromSlash(files[i]))
	}

	err = createGlobFiles(files)
	if err != nil {
		t.Fatal(err)
	}

	globs := []string{
		filepath.Join(dir, "**", "*.json"),
	}

	matches, err := getGlobFiles(globs)

	if err != nil {
		t.Error(err)
	}

	sort.Strings(files)
	sort.Strings(matches)

	if !reflect.DeepEqual(files, matches) {
		t.Errorf("Value %#v not equal expected %#v\n", files, matches)
	}
}

func TestSlashGlob(t *testing.T) {
	dir, err := createTempDir()

	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	files := []string{
		"test.json",
		"woot/test.json",
		"woot/abc/test.json",
		"gg/t/a.json",
	}

	for i := range files {
		files[i] = filepath.Join(dir, filepath.FromSlash(files[i]))
	}

	err = createGlobFiles(files)
	if err != nil {
		t.Fatal(err)
	}

	globs := []string{
		filepath.ToSlash(filepath.Join(dir, "**", "*.json")),
	}

	matches, err := getGlobFiles(globs)

	if err != nil {
		t.Error(err)
	}

	sort.Strings(files)
	sort.Strings(matches)

	if !reflect.DeepEqual(files, matches) {
		t.Errorf("Value %#v not equal expected %#v\n", files, matches)
	}
}
