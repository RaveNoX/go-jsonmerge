package options

import (
	"fmt"
	"os"

	"errors"

	"github.com/juju/gnuflag"
)

var (
	// ErrorVerboseAndQuiet indicates than user specified "verbose" and "quiet" flags together
	ErrorVerboseAndQuiet = errors.New(`You cannot use "verbose" and "quiet" flags together`)

	// ErrorPatchAndQuiet indicates than user specified "patch" and "quiet" flags together
	ErrorPatchAndQuiet = errors.New(`You cannot use "patch" and "quiet" flags together`)

	// ErrorVerboseAndPatch indicates than user specified "verbose" and "patch" flags together
	ErrorVerboseAndPatch = errors.New(`You cannot use "verbose" and "patch" flags together`)

	// ErrorNotEnoughArguments indicates than user passed too few arguments
	ErrorNotEnoughArguments = errors.New(`Not enough arguments`)
)

// Options for application
type Options struct {
	Verbose, Quiet, DryRun, Patch bool
	PatchFile                     string
	Globs                         []string
	Name                          string
}

func (options *Options) getFlags() *gnuflag.FlagSet {
	flags := gnuflag.NewFlagSet(options.Name, gnuflag.ContinueOnError)

	flags.BoolVar(&options.Patch, "patch", false, "Print patch")
	flags.BoolVar(&options.DryRun, "dry", false, "Dry run: do not really make changes")
	flags.BoolVar(&options.Quiet, "quiet", false, "Do not display changed files")
	flags.BoolVar(&options.Verbose, "verbose", false, "Display changed values")
	flags.BoolVar(&options.Patch, "p", false, "Print patch")
	flags.BoolVar(&options.DryRun, "d", false, "Dry run: do not really make changes")
	flags.BoolVar(&options.Quiet, "q", false, "Do not display changed files")
	flags.BoolVar(&options.Verbose, "v", false, "Display changed values")

	return flags
}

// Parse parses options, emits error if any
func (options *Options) Parse(arguments []string) error {
	flags := options.getFlags()

	err := flags.Parse(false, arguments)

	if err != nil {
		return err
	}

	args := flags.Args()

	if options.Verbose && options.Quiet {
		return ErrorVerboseAndQuiet

	}

	if options.Patch && options.Quiet {
		return ErrorPatchAndQuiet
	}

	if options.Verbose && options.Patch {
		return ErrorVerboseAndPatch
	}

	if len(args) < 2 {
		return ErrorNotEnoughArguments
	}

	options.PatchFile = args[0]
	options.Globs = args[1:]

	return nil
}

func (options *Options) printUsage() {
	flags := options.getFlags()

	format := "%s\n    %s\n"
	fmt.Fprintf(os.Stderr, "Usage: %s: [args] <patch> <glob1>..<globN>\n", options.Name)
	fmt.Fprintf(os.Stderr, format, "<patch>", "Path to patch json file, use \"-\" for STDIN")
	fmt.Fprintf(os.Stderr, format, "<glob>", "Double star glob")
	flags.PrintDefaults()
}

// ParseOrExit parses options, if any error calls os.Exit(2)
func (options *Options) ParseOrExit(arguments []string) {
	err := options.Parse(arguments)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		options.printUsage()
		os.Exit(2)
	}
}
