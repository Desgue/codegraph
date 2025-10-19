# Data Model: Parse Directory CLI

**Feature**: 001-parse-directory-cli | **Date**: 2025-10-19

## Overview

This document defines the domain entities and their relationships for the CLI command parsing and directory validation feature.

## Entities

### ParseCommand

Represents the `parse` CLI command with its arguments and flags.

**Responsibilities**:
- Parse and validate command-line flags
- Coordinate directory resolution and validation
- Provide user feedback (error messages, informational logs)

**Fields**:

| Field | Type | Description | Validation Rules |
|-------|------|-------------|------------------|
| `TargetDirectory` | `*TargetDirectory` | The directory to parse | Must not be nil after initialization |
| `OutputFile` | `string` | Path to output file (--output flag) | Must be non-empty string |
| `IncludeTests` | `bool` | Whether to include test files (--include-tests flag) | No validation (boolean) |

**Methods**:

```go
// NewParseCommand creates a new ParseCommand by parsing flags and validating inputs
// Returns error if flags are invalid or directory validation fails
func NewParseCommand(args []string) (*ParseCommand, error)

// Validate ensures all command fields are valid
// Called internally by NewParseCommand, exposed for testing
func (pc *ParseCommand) Validate() error

// Execute runs the parse command (placeholder for now, returns nil)
// Future implementation will contain actual parsing logic
func (pc *ParseCommand) Execute() error
```

**State Transitions**: None (command is created validated or returns error)

**Invariants**:
- If `ParseCommand` exists (not error), all fields are validated
- `TargetDirectory` is never nil if command was successfully created
- `OutputFile` is never empty string

---

### TargetDirectory

Represents a validated, resolved directory path.

**Responsibilities**:
- Validate directory exists and is accessible
- Resolve relative paths to absolute paths
- Follow symbolic links to actual directory
- Encapsulate all directory path validation logic

**Fields**:

| Field | Type | Description | Validation Rules |
|-------|------|-------------|------------------|
| `InputPath` | `string` | Original path provided by user | Preserved for error messages |
| `ResolvedPath` | `string` | Absolute path after resolution | Must be absolute, must exist |
| `WasDefaulted` | `bool` | True if current directory was used as default | Read-only after creation |

**Methods**:

```go
// NewTargetDirectory creates a validated TargetDirectory from user input
// inputPath: empty string means "use current directory"
// Returns error if path doesn't exist, isn't a directory, or lacks permissions
func NewTargetDirectory(inputPath string) (*TargetDirectory, error)

// Validate checks if the resolved path is a valid, accessible directory
// Returns specific errors for: not exists, is file, permission denied
func (td *TargetDirectory) Validate() error

// String returns the resolved absolute path (for display/logging)
func (td *TargetDirectory) String() string

// LogDefaultBehavior writes informational message to stderr if directory was defaulted
func (td *TargetDirectory) LogDefaultBehavior()
```

**State Transitions**: Immutable after creation (value object pattern)

**Invariants**:
- `ResolvedPath` is always absolute (no relative components like `.`, `..`)
- If `TargetDirectory` exists (not error), the path exists on filesystem and is a directory
- `InputPath` matches `ResolvedPath` unless path was relative or defaulted
- `WasDefaulted` is true only if `InputPath` was empty string

**Validation Rules**:

1. **Existence Check**:
   - Path must exist on filesystem
   - Error: `"directory does not exist: %s"`

2. **Type Check**:
   - Path must be a directory, not a file
   - Use `FileInfo.IsDir()` from `os.Stat`
   - Error: `"'%s' is a file, not a directory"`

3. **Permission Check**:
   - User must have read permission for directory
   - Detected via `os.IsPermission(err)` from `os.Stat` failure
   - Error: `"permission denied accessing '%s'"`

4. **Symbolic Link Handling**:
   - Automatically followed by `os.Stat` (not `os.Lstat`)
   - `ResolvedPath` is the actual directory, not the symlink
   - No explicit error for symlinks - handled transparently

---

## Relationships

```
ParseCommand "1" ──owns──> "1" TargetDirectory
    │
    └──> OutputFile (string)
    └──> IncludeTests (bool)
```

**Cardinality**:
- One `ParseCommand` has exactly one `TargetDirectory`
- `TargetDirectory` is owned by `ParseCommand` (composition, not aggregation)

**Lifecycle**:
1. User invokes CLI with args: `codegraph parse /path --output file.graphml`
2. `ParseCommand.NewParseCommand(args)` is called
3. Inside constructor:
   - Flags are parsed using `flag.NewFlagSet`
   - `TargetDirectory.NewTargetDirectory(dirArg)` is called
   - `TargetDirectory` validates path (exists, is dir, has permissions)
   - `ParseCommand` validates flags (output is non-empty)
4. If all validations pass, `ParseCommand` is returned
5. Caller can invoke `ParseCommand.Execute()` (currently no-op)

**Error Propagation**:
- `TargetDirectory` errors bubble up to `ParseCommand`
- `ParseCommand` wraps errors with command context
- All errors returned to main, written to `os.Stderr`, exit code 1

---

## Validation Matrix

| Scenario | InputPath | Expected Behavior | Error Message |
|----------|-----------|-------------------|---------------|
| Valid absolute path | `/home/user/project` | ✅ Resolve, validate | None |
| Valid relative path | `./subdirectory` | ✅ Resolve to absolute, validate | None |
| Current directory | `""` (empty) | ✅ Use `os.Getwd()`, log default | None (info log) |
| Nonexistent path | `/does/not/exist` | ❌ Error | `"directory does not exist: /does/not/exist"` |
| File path | `/home/user/file.go` | ❌ Error | `"'/home/user/file.go' is a file, not a directory"` |
| Permission denied | `/root/private` | ❌ Error | `"permission denied accessing '/root/private'"` |
| Symlink to directory | `/link -> /real/dir` | ✅ Follow, resolve to `/real/dir` | None |
| Symlink to file | `/link -> /file.txt` | ❌ Error | `"'[resolved path]' is a file, not a directory"` |
| Path with spaces | `/path with spaces/dir` | ✅ Handle correctly (Go does this) | None |
| Path with unicode | `/путь/目录` | ✅ Handle correctly (UTF-8 support) | None |

---

## Domain Rules

### Rule 1: No Anemic Models
Both `ParseCommand` and `TargetDirectory` encapsulate behavior with their data. Validation logic lives in the entities themselves, not in separate service classes.

### Rule 2: Fail Fast
Validation happens at creation time (`NewParseCommand`, `NewTargetDirectory`). If constructor succeeds, object is guaranteed valid.

### Rule 3: Immutability
`TargetDirectory` is immutable after creation (value object). No setters. This prevents invalid state.

### Rule 4: Explicit Error Types
Validation errors are concrete, user-friendly messages. Use `fmt.Errorf` with context. No generic "invalid input" errors.

### Rule 5: Single Responsibility
- `ParseCommand`: CLI interface, flag parsing, user interaction
- `TargetDirectory`: Path validation, resolution, filesystem interaction

---

## Testing Scenarios

### Unit Tests for TargetDirectory

1. **Valid directory (absolute path)**
   - Input: `/tmp/test-dir` (exists)
   - Expected: `ResolvedPath = "/tmp/test-dir"`, no error

2. **Valid directory (relative path)**
   - Input: `./testdata/sample`
   - Expected: `ResolvedPath = "[cwd]/testdata/sample"`, no error

3. **Current directory default**
   - Input: `""` (empty string)
   - Expected: `ResolvedPath = [cwd]`, `WasDefaulted = true`, no error

4. **Nonexistent directory**
   - Input: `/does/not/exist`
   - Expected: Error contains "directory does not exist"

5. **File instead of directory**
   - Input: `/tmp/file.txt` (create temp file)
   - Expected: Error contains "is a file, not a directory"

6. **Permission denied**
   - Input: `/root` (on Unix, if not running as root)
   - Expected: Error contains "permission denied"

7. **Symbolic link to directory**
   - Input: `/tmp/link -> /tmp/real-dir`
   - Expected: `ResolvedPath = "/tmp/real-dir"`, no error

### Unit Tests for ParseCommand

1. **Valid command with all flags**
   - Args: `["parse", "/tmp/dir", "--output", "out.graphml", "--include-tests"]`
   - Expected: All fields populated correctly

2. **Valid command with minimal flags**
   - Args: `["parse", "/tmp/dir", "--output", "out.graphml"]`
   - Expected: `IncludeTests = false`

3. **Missing output flag**
   - Args: `["parse", "/tmp/dir"]`
   - Expected: Error contains "--output flag requires a file path"

4. **Empty output flag**
   - Args: `["parse", "/tmp/dir", "--output", ""]`
   - Expected: Error contains "--output flag requires a file path"

5. **Invalid directory**
   - Args: `["parse", "/does/not/exist", "--output", "out.graphml"]`
   - Expected: Error from TargetDirectory validation

6. **No directory argument (default to current)**
   - Args: `["parse", "--output", "out.graphml"]`
   - Expected: `TargetDirectory.WasDefaulted = true`, stderr log message

---

## Implementation Notes

### Package Organization

```
codegraph/
├── main.go                     # Entry point, routes to commands
├── cli/
│   ├── parse_command.go        # ParseCommand entity
│   └── parse_command_test.go   # Table-driven tests
├── path/
│   ├── validator.go            # TargetDirectory entity
│   └── validator_test.go       # Table-driven tests with temp dirs
```

**Package Naming Rationale**:
- `cli`: Short, clear, combines well with types (`cli.ParseCommand`)
- `path`: Domain-specific (not `util` or `helper`), describes responsibility

### Constructor Pattern

Both entities use the "New + Validate" pattern:
```go
func NewEntity(input) (*Entity, error) {
    entity := &Entity{...}
    if err := entity.Validate(); err != nil {
        return nil, err
    }
    return entity, nil
}
```

This ensures entities are never in invalid state ("make invalid states unrepresentable").

### Error Handling Convention

All validation errors:
- Are returned, not panicked
- Include context (what path, what went wrong)
- Are user-facing (clear, actionable messages)
- Use `fmt.Errorf` for wrapping (`%w` for error chains)

Example:
```go
if err != nil {
    return fmt.Errorf("failed to validate directory '%s': %w", path, err)
}
```
