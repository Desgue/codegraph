# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Build the binary
make build

# Run tests
go test ./...

# Run a specific test
go test -run TestName ./path/to/package
```

## Project Architecture

**codegraph** is a CLI tool for parsing and analyzing Go codebases to generate dependency graphs.

### Core Structure

- **main.go**: Entry point with subcommand routing. Currently supports the `parse` command.
- **cli/**: Command implementations
  - `ParseCommand`: Handles the `parse` subcommand with flag parsing and validation
- **path/**: Path resolution and validation
  - `TargetDirectory`: Validates and resolves directory paths, handling symlinks and permission checks

### Command Flow

1. `main.go` routes to the appropriate subcommand handler
2. Command constructors (e.g., `NewParseCommand`) parse flags and validate inputs
3. `path.NewTargetDirectory` resolves paths:
   - Defaults to current working directory if no path provided
   - Converts relative paths to absolute paths
   - Resolves symlinks to canonical paths using `filepath.EvalSymlinks`
   - Validates directory exists and is accessible
4. Command `Execute()` methods perform the actual work

### Path Resolution Rules

- Empty input → current working directory
- Relative paths → converted to absolute
- Symlinks → resolved to canonical paths
- Files rejected with error message
- Permission errors fail immediately

### Current Feature Branch

Working on `001-parse-directory-cli`: Implementing the basic CLI foundation with directory path handling and validation. See `specs/001-parse-directory-cli/` for detailed requirements.
