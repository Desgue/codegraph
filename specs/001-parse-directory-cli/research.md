# Research: Parse Directory CLI

**Feature**: 001-parse-directory-cli | **Date**: 2025-10-19

## Overview

This document captures research findings for implementing a CLI tool using Go's standard library `flag` package with directory validation capabilities.

## Decision 1: CLI Command Pattern

**Decision**: Use subcommand pattern with `flag.NewFlagSet` for the `parse` command

**Rationale**:
- Go's `flag` package doesn't have built-in subcommand support like cobra/cli libraries
- `flag.NewFlagSet` allows creating independent flag sets for each subcommand
- Enables future extensibility (adding more commands like `export`, `analyze`) without external dependencies
- Keeps implementation simple and idiomatic to Go stdlib

**Alternatives Considered**:
1. Single global FlagSet - rejected because it doesn't scale to multiple commands
2. External library (cobra, cli) - rejected due to "standard library only" constraint
3. Manual argument parsing - rejected because `flag` provides validation, help text generation

**Implementation Pattern**:
```go
// main.go
func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Usage: codegraph <command> [options]\n")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "parse":
        parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)
        // register flags here
        parseCmd.Parse(os.Args[2:])
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}
```

## Decision 2: Directory Validation Strategy

**Decision**: Use `os.Stat` + `os.IsNotExist` + `FileInfo.IsDir()` for path validation

**Rationale**:
- `os.Stat` provides file metadata including directory status, permissions info
- Follows symbolic links by default (spec requirement FR-009)
- `os.IsNotExist(err)` explicitly differentiates "not found" from permission errors
- Single syscall, minimal overhead (<1ms typically)

**Alternatives Considered**:
1. `os.Lstat` - rejected because spec requires following symlinks, Lstat doesn't follow
2. `filepath.Walk` - rejected as overkill for just checking if path is directory
3. Custom syscall - rejected as unnecessary complexity

**Error Handling Pattern**:
```go
func ValidateDirectory(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("directory does not exist: %s", path)
        }
        if os.IsPermission(err) {
            return fmt.Errorf("permission denied accessing '%s'", path)
        }
        return fmt.Errorf("error accessing '%s': %w", path, err)
    }

    if !info.IsDir() {
        return fmt.Errorf("'%s' is a file, not a directory", path)
    }

    return nil
}
```

## Decision 3: Path Resolution Approach

**Decision**: Use `filepath.Abs` for resolving relative paths to absolute paths

**Rationale**:
- Spec requires resolving relative paths (FR-008)
- `filepath.Abs` handles `.`, `..`, and relative paths correctly across platforms
- Automatically prepends current working directory
- Cross-platform (handles Windows `C:\` and Unix `/` correctly)
- Returns cleaned path (removes redundant separators, resolves `.` and `..`)

**Alternatives Considered**:
1. `filepath.Clean` only - rejected because it doesn't make relative paths absolute
2. Manual path joining with `os.Getwd()` - rejected as reinventing stdlib functionality
3. `filepath.EvalSymlinks` - rejected because it fails if path doesn't exist yet

**Pattern**:
```go
func ResolvePath(inputPath string) (string, error) {
    absolutePath, err := filepath.Abs(inputPath)
    if err != nil {
        return "", fmt.Errorf("failed to resolve path '%s': %w", inputPath, err)
    }
    return absolutePath, nil
}
```

## Decision 4: Current Directory Defaulting

**Decision**: Use `os.Getwd()` when no directory argument provided

**Rationale**:
- Spec requirement FR-006: use current working directory as default
- `os.Getwd()` is the standard Go idiom for getting current directory
- Returns absolute path, no need for additional resolution
- Fails explicitly if current directory was deleted or has permission issues

**Alternatives Considered**:
1. `filepath.Abs(".")` - works but adds unnecessary step (Abs calls Getwd internally)
2. Environment variable `PWD` - rejected as unreliable (can be stale/wrong)

**Pattern**:
```go
func GetTargetDirectory(args []string) (string, error) {
    if len(args) == 0 {
        cwd, err := os.Getwd()
        if err != nil {
            return "", fmt.Errorf("failed to get current directory: %w", err)
        }
        fmt.Fprintf(os.Stderr, "No directory specified, using current directory: %s\n", cwd)
        return cwd, nil
    }
    return args[0], nil
}
```

## Decision 5: Flag Validation for --output

**Decision**: Validate `--output` is non-empty string, defer file creation until actual output

**Rationale**:
- User specified "validation for flag inputs for now" - no actual file I/O yet
- Empty string is invalid output path (would cause error later)
- Don't validate parent directory exists yet - that's implementation detail for later
- Keep validation minimal per YAGNI principle

**Pattern**:
```go
outputFile := parseCmd.String("output", "", "Output file path (required)")
parseCmd.Parse(os.Args[2:])

if *outputFile == "" {
    return fmt.Errorf("--output flag requires a file path")
}
```

## Decision 6: Flag Validation for --include-tests

**Decision**: Boolean flag, no validation needed beyond type check

**Rationale**:
- Boolean flags are self-validating (true/false only)
- Go's `flag.Bool` handles parsing automatically
- No edge cases to validate

**Pattern**:
```go
includeTests := parseCmd.Bool("include-tests", false, "Include test files in parsing")
// No validation needed - type-safe by design
```

## Best Practices for Go CLI with stdlib flag

### Error Handling
- Write errors to `os.Stderr`, not `os.Stdout`
- Use `fmt.Fprintf(os.Stderr, ...)` for error messages
- Exit with code 1 for user errors, code 2 for internal errors
- Wrap errors with context using `fmt.Errorf("context: %w", err)`

### Help Text
- Let `flag` package generate help automatically via `flag.ExitOnError`
- Provide clear, concise descriptions for each flag
- Include usage example in error messages when command is missing

### Logging vs Output
- Logs (informational): `os.Stderr`
- Actual output data: `os.Stdout`
- Rationale: Allows piping output without log noise: `codegraph parse | tool`

### Testing
- Use `flag.NewFlagSet` instead of global `flag` package in production code
- Easier to test - can create isolated flagsets per test
- Avoid global state contamination between tests

## Cross-Platform Considerations

### Path Separators
- Always use `filepath.Join()` or `path/filepath` functions, never hardcode `/` or `\`
- `filepath.Abs`, `filepath.Clean` handle platform differences automatically

### Case Sensitivity
- Filesystem paths are case-sensitive on Linux/macOS, case-insensitive on Windows
- Don't rely on case for validation - let OS handle it via `os.Stat`

### Line Endings
- Use `\n` in Go code - Go runtime handles conversion on Windows
- `fmt.Fprintf` and `fmt.Println` work correctly cross-platform

## Performance Notes

- `os.Stat`: Typically <1ms for local filesystem, can be slower on network mounts
- `filepath.Abs`: ~100μs (microseconds), negligible overhead
- `os.Getwd`: ~50μs, cached by OS, very fast
- Flag parsing: <100μs for small flag sets (2-3 flags)

**Total validation overhead**: <10ms target easily achievable with this approach

## Security Considerations

### Path Traversal
- `filepath.Clean` and `filepath.Abs` automatically prevent `../../../` attacks
- Always resolve to absolute path before validation
- Spec follows symlinks - this is intentional, not a security issue (user must have permissions)

### Permission Checking
- `os.Stat` respects OS permissions - if user can't read directory, validation fails appropriately
- No need to manually check Unix permission bits - OS does this

## References

- Go stdlib `flag` documentation: https://pkg.go.dev/flag
- Go stdlib `os` documentation: https://pkg.go.dev/os
- Go stdlib `path/filepath` documentation: https://pkg.go.dev/path/filepath
- Effective Go - Command-line flags: https://go.dev/doc/effective_go#flags
