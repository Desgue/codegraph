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
  - `ParseCommand`: Handles the `parse` subcommand with `--output` and `--include-tests` flags
- **path/**: Path resolution and validation
  - `TargetDirectory`: Validates and resolves directory paths, handling symlinks and permission checks
- **parser/**: Go source code parsing using `golang.org/x/tools/go/packages`
  - `Load()`: Parses packages with automatic deduplication and test variant handling
  - Returns AST with syntax trees, imports, and type information

### Command Flow

1. `main.go` routes to the appropriate subcommand handler
2. Command constructors (e.g., `NewParseCommand`) parse flags and validate inputs
3. `path.NewTargetDirectory` resolves paths:
   - Defaults to current working directory if no path provided
   - Converts relative paths to absolute paths
   - Resolves symlinks to canonical paths using `filepath.EvalSymlinks`
   - Validates directory exists and is accessible
4. `parser.Load()` uses `go/packages` to parse Go code:
   - Loads with `NeedName`, `NeedFiles`, `NeedSyntax`, `NeedImports`, `NeedTypes` modes
   - Automatically deduplicates package variants (when `includeTests=true`, keeps variant with most files)
   - Filters out synthetic `.test` packages
   - Returns deterministically ordered packages by import path
5. Command `Execute()` methods output parsed results

### Path Resolution Rules

- Empty input → current working directory
- Relative paths → converted to absolute
- Symlinks → resolved to canonical paths
- Files rejected with error message
- Permission errors fail immediately

### Parser Behavior

- **Test Handling**: When `--include-tests=true`, deduplicates package variants by keeping the one with the most files (test variant includes both production and test files)
- **Error Handling**: Package-level parse errors are counted and reported via `packages.PrintErrors()`, but don't fail the entire operation
- **Multi-module Limitation**: Module path detection uses the first discovered module; monorepos with multiple modules are not fully supported
- **Comment Preservation**: AST includes comments via `NeedSyntax` for future documentation analysis
