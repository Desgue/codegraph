package main

import (
	"fmt"
	"os"

	"github.com/Desgue/codegraph/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: codegraph <command> [options]\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "parse":
		parseCommand, err := cli.NewParseCommand(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := parseCommand.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nUsage: codegraph <command> [options]\n", os.Args[1])
		os.Exit(1)
	}
}
