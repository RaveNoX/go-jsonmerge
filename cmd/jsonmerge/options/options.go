package options

import (
	"fmt"
	"os"

	"github.com/juju/gnuflag"
)

// Options for application
type Options struct {
	Verbose, Quiet bool
	Patch          string
	Globs          []string
	Name           string
}

func (options *Options) getFlags() *gnuflag.FlagSet {
	flags := gnuflag.NewFlagSet(options.Name, gnuflag.ContinueOnError)

	flags.BoolVar(&options.Quiet, "quiet", false, "Do not display changed files")
	flags.BoolVar(&options.Verbose, "verbose", false, "Display changed values")
	flags.BoolVar(&options.Quiet, "q", false, "Do not display changed files")
	flags.BoolVar(&options.Verbose, "v", false, "Display changed values")

	return flags
}

// Parse parses options, emits error if any
func (options *Options) Parse(arguments []string) (err error) {
	flags := options.getFlags()

	err = flags.Parse(false, arguments)

	if err != nil {
		return
	}

	args := flags.Args()

	if options.Verbose && options.Quiet {
		err = fmt.Errorf("You cannot use \"verbose\" and \"quiet\" flags together")
		return
	}

	if len(args) < 2 {
		err = fmt.Errorf("Not enough arguments")
		return
	}

	options.Patch = args[0]
	options.Globs = args[1:]

	return
}

func (options *Options) printUsage() {
	flags := options.getFlags()

	format := "%s\n    %s\n"
	fmt.Fprintf(os.Stderr, "Usage: %s: [args] <patch> <glob1>..<globN>\n", options.Name)
	fmt.Fprintf(os.Stderr, format, "<patch>", "Path to patch json file")
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
