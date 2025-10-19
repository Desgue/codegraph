package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/Desgue/codegraph/parser"
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
	includeTests := flagSet.Bool("include-tests", true, "Include test files in parsing")

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
	pkgs, errorCount, err := parser.Load(pc.TargetDirectory.Path, pc.IncludeTests)
	if err != nil {
		return err
	}

	totalPackages := len(pkgs)
	totalFiles := 0
	var modulePath string

	for _, pkg := range pkgs {
		fmt.Printf("\nPackage: %s\n", pkg.PkgPath)
		fmt.Printf("  Name: %s\n", pkg.Name)
		fmt.Printf("  Files (%d):\n", len(pkg.GoFiles))
		for _, file := range pkg.GoFiles {
			fmt.Printf("    - %s\n", file)
		}
		if len(pkg.Errors) > 0 {
			fmt.Printf("  Errors: %d\n", len(pkg.Errors))
		}

		totalFiles += len(pkg.GoFiles)
		// Module path detection assumes all packages belong to the same Go module.
		// Uses the first non-nil Module found.
		// LIMITATION: Multi-module repositories (monorepos) are not supported.
		// Only the first discovered module path will be displayed in the summary.
		if pkg.Module != nil && modulePath == "" {
			modulePath = pkg.Module.Path
		}
	}

	fmt.Printf("\n")
	if modulePath != "" {
		fmt.Printf("Module: %s\n", modulePath)
	}
	fmt.Printf("Loaded %d packages, parsed %d files\n", totalPackages, totalFiles)
	if errorCount > 0 {
		fmt.Fprintf(os.Stderr, "Encountered %d parse errors\n", errorCount)
	}

	return nil
}
