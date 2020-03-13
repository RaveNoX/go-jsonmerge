package options

import (
	"testing"
)

func TestArgsLen(t *testing.T) {
	var args []string

	opts := new(Options)

	err := opts.Parse(args)
	if err != ErrorNotEnoughArguments {
		t.Errorf("Expected %v error got %v\n", ErrorNotEnoughArguments, err)
	}

	args = append(args, "patch.json")

	err = opts.Parse(args)
	if err != ErrorNotEnoughArguments {
		t.Errorf("Expected %v error got %v\n", ErrorNotEnoughArguments, err)
	}

	args = append(args, "glob1.json")

	err = opts.Parse(args)
	if err != nil {
		t.Errorf("Expected nil error got %v\n", err)
	}

	args = append(args, "glob2.json")

	err = opts.Parse(args)
	if err != nil {
		t.Errorf("Expected nil error got %v\n", err)
	}
}

func TestPatch(t *testing.T) {
	args := []string{
		"--patch",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.Patch {
		t.Error("Patch not set")
	}
}

func TestPatchShort(t *testing.T) {
	args := []string{
		"-p",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.Patch {
		t.Error("Patch not set")
	}
}

func TestVerbose(t *testing.T) {
	args := []string{
		"--verbose",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.Verbose {
		t.Error("Verbose not set")
	}
}

func TestVerboseShort(t *testing.T) {
	args := []string{
		"-v",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.Verbose {
		t.Error("Verbose not set")
	}
}

func TestDry(t *testing.T) {
	args := []string{
		"--dry",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.DryRun {
		t.Error("Dry not set")
	}
}

func TestDryShort(t *testing.T) {
	args := []string{
		"-d",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.DryRun {
		t.Error("Dry not set")
	}
}

func TestQuiet(t *testing.T) {
	args := []string{
		"--quiet",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.Quiet {
		t.Error("Quiet not set")
	}
}

func TestQuietShort(t *testing.T) {
	args := []string{
		"-q",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != nil {
		t.Error(err)
	}

	if !opts.Quiet {
		t.Error("Quiet not set")
	}
}

func TestVerboseAndQuiet(t *testing.T) {
	args := []string{
		"-q",
		"-v",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != ErrorVerboseAndQuiet {
		t.Errorf("Expected %v error got %v\n", ErrorVerboseAndQuiet, err)
	}

	args = []string{
		"--quiet",
		"-v",
		"path.json",
		"glob.json",
	}

	err = opts.Parse(args)

	if err != ErrorVerboseAndQuiet {
		t.Errorf("Expected %v error got %v\n", ErrorVerboseAndQuiet, err)
	}

	args = []string{
		"--quiet",
		"--verbose",
		"path.json",
		"glob.json",
	}

	err = opts.Parse(args)

	if err != ErrorVerboseAndQuiet {
		t.Errorf("Expected %v error got %v\n", ErrorVerboseAndQuiet, err)
	}
}

func TestPatchAndQuiet(t *testing.T) {
	args := []string{
		"-q",
		"-p",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != ErrorPatchAndQuiet {
		t.Errorf("Expected %v error got %v\n", ErrorPatchAndQuiet, err)
	}

	args = []string{
		"--quiet",
		"-p",
		"path.json",
		"glob.json",
	}

	err = opts.Parse(args)

	if err != ErrorPatchAndQuiet {
		t.Errorf("Expected %v error got %v\n", ErrorPatchAndQuiet, err)
	}

	args = []string{
		"--quiet",
		"--patch",
		"path.json",
		"glob.json",
	}

	err = opts.Parse(args)

	if err != ErrorPatchAndQuiet {
		t.Errorf("Expected %v error got %v\n", ErrorPatchAndQuiet, err)
	}
}

func TestVerboseAndPatch(t *testing.T) {
	args := []string{
		"-v",
		"-p",
		"path.json",
		"glob.json",
	}

	opts := new(Options)
	err := opts.Parse(args)

	if err != ErrorVerboseAndPatch {
		t.Errorf("Expected %v error got %v\n", ErrorVerboseAndPatch, err)
	}

	args = []string{
		"--verbose",
		"-p",
		"path.json",
		"glob.json",
	}

	err = opts.Parse(args)

	if err != ErrorVerboseAndPatch {
		t.Errorf("Expected %v error got %v\n", ErrorVerboseAndPatch, err)
	}

	args = []string{
		"--verbose",
		"--patch",
		"path.json",
		"glob.json",
	}

	err = opts.Parse(args)

	if err != ErrorVerboseAndPatch {
		t.Errorf("Expected %v error got %v\n", ErrorVerboseAndPatch, err)
	}
}
