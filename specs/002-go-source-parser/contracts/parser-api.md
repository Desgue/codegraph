# Parser Package API Contract

**Package**: `github.com/Desgue/codegraph/parser`
**Version**: 1.0.0
**Date**: 2025-10-19

## Overview

The parser package provides Go source code parsing functionality using `golang.org/x/tools/go/packages`. It discovers packages, parses files, and aggregates results for downstream graph construction.

## Public API

### Load Function

```go
func Load(targetDir string) (*Result, error)
```

**Description**: Parses all Go packages in the target directory and returns aggregated results.

**Parameters**:
- `targetDir` (string): Absolute path to directory containing Go source code
  - Must be a valid directory path
  - Must have read permissions
  - Should be pre-validated by `path.TargetDirectory`

**Returns**:
- `*Result`: Parsing results with statistics and package data
- `error`: Non-nil only for catastrophic failures (pattern parsing, driver issues)

**Behavior**:
- Discovers all packages matching pattern `"./..."` relative to targetDir
- Parses files with LoadMode: `NeedName | NeedFiles | NeedSyntax | NeedImports | NeedTypes`
- Reports parse errors to stderr via `packages.PrintErrors()`
- Returns partial results when some packages/files fail to parse
- Only returns error for catastrophic failures (not individual parse errors)

**Error Cases**:
- Returns error: Invalid pattern, driver initialization failure, nil targetDir
- Returns nil error: Parse errors, missing imports, syntax errors (captured in Result.ErrorCount)

**Performance**:
- Target: <5 seconds for 50,000 LOC
- Supports up to 1000 source files

**Example**:
```go
result, err := parser.Load("/path/to/project")
if err != nil {
    // Catastrophic failure - exit 1
    return fmt.Errorf("failed to parse: %w", err)
}

// Partial success - exit 0
fmt.Printf("Packages: %d, Files: %d, Errors: %d\n",
    result.TotalPackages, result.TotalFiles, result.ErrorCount)
```

### Result Type

```go
type Result struct {
    Packages      []*packages.Package
    TotalPackages int
    TotalFiles    int
    ErrorCount    int
    SkippedDirs   []string
}
```

**Description**: Aggregates parsing statistics and package data.

**Fields**:

#### Packages
- **Type**: `[]*packages.Package`
- **Description**: All discovered packages with parsed AST data
- **Contains**: Package name, import path, file paths, syntax trees, imports, errors, module info
- **Usage**: Pass to graph construction phase for dependency analysis

#### TotalPackages
- **Type**: `int`
- **Description**: Count of packages discovered
- **Range**: `>= 0`
- **Semantics**: `len(Packages)`

#### TotalFiles
- **Type**: `int`
- **Description**: Count of .go files processed
- **Range**: `>= 0`
- **Semantics**: Sum of `len(pkg.GoFiles)` across all packages

#### ErrorCount
- **Type**: `int`
- **Description**: Count of files that failed to parse
- **Range**: `>= 0`
- **Semantics**: Return value from `packages.PrintErrors()`

#### SkippedDirs
- **Type**: `[]string`
- **Description**: Directories skipped due to permission errors
- **Format**: Absolute directory paths
- **Empty**: When all directories are accessible

**Invariants**:
- `len(Packages) == TotalPackages`
- `TotalFiles >= TotalPackages` (at least one file per package)
- `ErrorCount >= 0`

**Usage**:
```go
for _, pkg := range result.Packages {
    fmt.Printf("Package: %s (%d files)\n", pkg.PkgPath, len(pkg.GoFiles))
    for _, err := range pkg.Errors {
        fmt.Printf("  Error: %s\n", err)
    }
}
```

## Dependency Contracts

### Input Contract (path package)

The parser expects input from `path.TargetDirectory`:

```go
targetDir, err := path.NewTargetDirectory(userInput)
if err != nil {
    return err
}
result, err := parser.Load(targetDir.Path())
```

**Assumptions**:
- `targetDir` is an absolute path
- `targetDir` exists and is a directory
- `targetDir` has read permissions

### Output Contract (future graph package)

The parser provides output for graph construction:

```go
result, err := parser.Load(targetDir)
if err != nil {
    return err
}
// Pass packages to graph builder
graph := builder.Build(result.Packages)
```

**Guarantees**:
- Each package has parsed syntax trees in `Syntax` field
- Import relationships are available in `Imports` field
- Module information is available in `Module` field (if in a module)

## External Dependencies

### golang.org/x/tools/go/packages

**Used For**: Package loading and AST parsing

**Key Types Used**:
- `packages.Package`: Parsed package with AST data
- `packages.Config`: Load configuration (Mode, Dir)
- `packages.Error`: Parse and type errors
- `packages.Module`: Module metadata

**Functions Used**:
- `packages.Load(cfg, patterns...)`: Load and parse packages
- `packages.PrintErrors(pkgs)`: Report errors to stderr

**Version Constraint**: Compatible with Go 1.24+

## Error Handling Contract

### Catastrophic Errors (exit 1)

Returned as `error` from `Load()`:
- Invalid pattern for `packages.Load`
- Driver initialization failure
- Nil or empty targetDir

**Client Action**: Return error, exit non-zero

### Partial Failures (exit 0)

Captured in `Result.ErrorCount`:
- Parse errors in Go files
- Missing imports
- Type-checking failures
- Permission errors on subdirectories

**Client Action**: Print statistics, exit zero

### Error Output Format

Errors printed to stderr by `packages.PrintErrors()`:
```
path/to/file.go:15:3: expected ';', found 'EOF'
path/to/other.go:42:1: undefined: unknownFunc
```

**Format**: `<file>:<line>:<col>: <message>`

## Thread Safety

The `Load()` function is **not thread-safe**:
- Concurrent calls with different targetDir: Safe
- Concurrent calls with same targetDir: Undefined behavior

The `Result` type is **immutable** after construction:
- Safe to read from multiple goroutines
- Do not modify Packages slice or individual packages

## Testing Contracts

### Unit Tests

```go
func TestLoad_ValidProject(t *testing.T)
func TestLoad_EmptyDirectory(t *testing.T)
func TestLoad_WithSyntaxErrors(t *testing.T)
func TestLoad_WithPermissionErrors(t *testing.T)
```

### Test Fixtures

Expected structure:
```
testdata/
├── valid-project/
│   └── main.go
├── empty-dir/
├── syntax-errors/
│   └── invalid.go
└── multi-package/
    ├── pkg1/
    └── pkg2/
```

### Mock Requirements

No mocks needed - test with real filesystem fixtures

## Versioning

**Semantic Versioning**: Major.Minor.Patch

**Current Version**: 1.0.0

**Breaking Changes** (require major version bump):
- Changing `Load()` signature
- Removing/renaming `Result` fields
- Changing error return semantics

**Non-Breaking Changes** (minor/patch):
- Adding new fields to `Result`
- Performance improvements
- Bug fixes

## Future API Extensions

Potential additions (not committed):

```go
// Load with context for cancellation
func LoadWithContext(ctx context.Context, targetDir string) (*Result, error)

// Load with custom configuration
func LoadWithConfig(targetDir string, cfg *Config) (*Result, error)

// Incremental parse (watch mode)
func Watch(targetDir string, onChange func(*Result)) error
```

## CLI Integration Contract

```go
// In cli/parse.go
func (c *ParseCommand) Execute() error {
    // 1. Validate directory
    targetDir, err := path.NewTargetDirectory(c.targetPath)
    if err != nil {
        return err
    }

    // 2. Parse packages
    result, err := parser.Load(targetDir.Path())
    if err != nil {
        return err  // Catastrophic error - exit 1
    }

    // 3. Print statistics - exit 0 even with parse errors
    fmt.Printf("Packages found: %d\n", result.TotalPackages)
    fmt.Printf("Files parsed: %d\n", result.TotalFiles)
    if result.ErrorCount > 0 {
        fmt.Printf("Parse errors: %d\n", result.ErrorCount)
    }
    if len(result.SkippedDirs) > 0 {
        fmt.Printf("Skipped directories: %d\n", len(result.SkippedDirs))
        for _, dir := range result.SkippedDirs {
            fmt.Printf("  - %s\n", dir)
        }
    }

    // 4. TODO: Pass result.Packages to graph construction
    return nil
}
```

## Performance Contract

**Guarantees**:
- Parse 50,000 LOC in <5 seconds (standard hardware)
- Support up to 1000 source files without degradation
- Memory usage proportional to AST size (no leak)

**Best Practices for Clients**:
- Call `Load()` once per directory, reuse Result
- Don't modify returned packages or AST nodes
- For large projects, consider progress reporting (future)

## Compatibility

**Go Version**: 1.24+
**Platforms**: Linux, macOS, Windows
**Module Mode**: Required (uses go.mod for resolution)
**GOPATH Mode**: Supported (fallback when no go.mod)
