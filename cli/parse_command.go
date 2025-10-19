package cli

import (
	"flag"
	"fmt"

	"github.com/Desgue/codegraph/path"
)

type ParseCommand struct {
	TargetDirectory *path.TargetDirectory
	OutputFile      string
	IncludeTests    bool
}

func NewParseCommand(args []string) (*ParseCommand, error) {
	flagSet := flag.NewFlagSet("parse", flag.ContinueOnError)

	outputFile := flagSet.String("output", "", "Output file path (required)")
	includeTests := flagSet.Bool("include-tests", false, "Include test files in parsing")

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	directoryArgument := ""
	if flagSet.NArg() > 0 {
		directoryArgument = flagSet.Arg(0)
	}

	targetDirectory, err := path.NewTargetDirectory(directoryArgument)
	if err != nil {
		return nil, err
	}

	parseCommand := &ParseCommand{
		TargetDirectory: targetDirectory,
		OutputFile:      *outputFile,
		IncludeTests:    *includeTests,
	}

	if err := parseCommand.Validate(); err != nil {
		return nil, err
	}

	return parseCommand, nil
}

func (pc *ParseCommand) Validate() error {
	if pc.OutputFile == "" {
		return fmt.Errorf("--output flag requires a file path")
	}
	return nil
}

func (pc *ParseCommand) Execute() error {
	return nil
}
