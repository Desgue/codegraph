# Data Model: Go Source Parser

**Date**: 2025-10-19
**Package**: parser

## Entities

### 1. packages.Package (External Type)

Represents a Go package with parsed AST data. This is the standard type from `golang.org/x/tools/go/packages`.

**Key Fields Used**:
- `Name`: `string` - Package name (e.g., "main", "parser")
- `PkgPath`: `string` - Import path (e.g., "github.com/Desgue/codegraph/parser")
- `GoFiles`: `[]string` - Absolute paths to .go files in package
- `Syntax`: `[]*ast.File` - Parsed AST for each file
- `Imports`: `map[string]*packages.Package` - Direct dependencies
- `Errors`: `[]packages.Error` - Parse/type errors for this package
- `Module`: `*packages.Module` - Module information if in a module

**Responsibilities**:
- Provide parsed AST structure for code analysis
- Track package-level errors
- Map import paths to package objects

**Design Decision**: No wrapper type - use `packages.Package` directly per YAGNI principle. Return `[]*packages.Package` directly from Load() function.

### 2. packages.Error (External Type)

Records parsing and type-checking errors.

**Key Fields**:
- `Pos`: `string` - File position in format "path/file.go:line:col"
- `Msg`: `string` - Error description
- `Kind`: `packages.ErrorKind` - Error category

**Responsibilities**:
- Record file path, line, column for each error
- Provide formatted error messages for terminal output

**Usage**: Printed via `packages.PrintErrors()` to stderr

## Relationships

```
[]*packages.Package (returned by Load())
    │
    ├── Name: string
    ├── PkgPath: string
    ├── GoFiles: []string
    ├── Syntax: []*ast.File
    ├── Imports: map[string]*packages.Package
    ├── Errors: []packages.Error
    │   ├── Pos: string
    │   ├── Msg: string
    │   └── Kind: ErrorKind
    └── Module: *packages.Module
```

## Package Structure

```
parser/
├── loader.go           # Load() function and Config setup
└── loader_test.go      # Load function tests
```

## Type Definitions

### loader.go

```go
package parser

import (
    "golang.org/x/tools/go/packages"
)

// Load parses all Go packages in targetDir and returns them.
// Returns error only for catastrophic failures (pattern parsing, driver issues).
// Package-level parse errors are printed to stderr via packages.PrintErrors().
// Each returned package contains Errors field with parse failures.
func Load(targetDir string) ([]*packages.Package, error)
```

## Domain Logic

### Parser Package Responsibilities

1. **Load(targetDir string) ([]*packages.Package, error)**
   - Configure packages.Load with LoadMode: NeedName | NeedFiles | NeedSyntax | NeedImports
   - Execute packages.Load with "./..." pattern
   - Call packages.PrintErrors to report errors to stderr
   - Return packages directly for CLI to process
   - Return error only for catastrophic failures

### Error Handling Strategy

- **Catastrophic errors**: Return from Load(), cause exit 1
  - Pattern parsing failures
  - Driver initialization failures
  - Directory not accessible

- **Package-level errors**: Include in Result, exit 0
  - Syntax errors in Go files
  - Missing imports
  - Type-checking failures

### Integration with CLI

```go
// In cli/parse.go Execute()
pkgs, err := parser.Load(targetDir.Path())
if err != nil {
    // Catastrophic error - exit 1
    return err
}

// Calculate statistics from packages
totalPackages := len(pkgs)
totalFiles := 0
errorCount := 0
for _, pkg := range pkgs {
    totalFiles += len(pkg.GoFiles)
    errorCount += len(pkg.Errors)
}

// Print statistics - exit 0 even with parse errors
fmt.Printf("Loaded %d packages, parsed %d files\n", totalPackages, totalFiles)
if errorCount > 0 {
    fmt.Fprintf(os.Stderr, "Encountered %d parse errors\n", errorCount)
}
```

## Future Extensions

The returned packages support future features:

1. **Graph Construction**: `pkgs[i].Imports` provides dependency relationships
2. **Documentation Analysis**: `pkgs[i].Syntax[j].Doc` contains comments
3. **Symbol Analysis**: `pkgs[i].Types` (if NeedTypes added later) provides type information
4. **Module Analysis**: `pkgs[i].Module` provides go.mod metadata

## Validation Rules

### Input Validation
- `targetDir` must be non-empty (validated by path.TargetDirectory)
- `targetDir` must exist and be a directory (validated by path.TargetDirectory)
- `targetDir` must have read permissions (validated by path.TargetDirectory)

### Output Validation
- `len(pkgs) >= 0` (may be 0 for empty directories)
- Each `pkg.GoFiles` contains valid file paths
- Each `pkg.Syntax` contains parsed AST trees
- `pkg.Errors` contains parse errors (if any)

## Error Scenarios

| Scenario | Handling | Exit Code |
|----------|----------|-----------|
| Directory contains no Go files | Return empty slice `[]` | 0 |
| Some files have syntax errors | Print errors to stderr, include in pkg.Errors | 0 |
| Permission denied on subdirectory | packages.Load skips, continues | 0 |
| Target directory doesn't exist | Return error from path.TargetDirectory | 1 |
| Invalid pattern for packages.Load | Return error from Load() | 1 |
| packages.Load driver failure | Return error from Load() | 1 |

## Performance Characteristics

- **Time Complexity**: O(n) where n is the number of files (linear scan)
- **Space Complexity**: O(m) where m is the total AST size (all files loaded in memory)
- **Target Performance**: 5 seconds for 50,000 LOC (100-200 files)
- **Memory Usage**: Moderate - full AST in memory but no type-checking data

## Design Decisions

### Why Return []*packages.Package Directly?

Per project constitution (KISS/YAGNI):
- No premature abstraction (no Result wrapper)
- Direct access to all package fields
- Future-proof (new fields automatically available)
- Standard documentation applies
- Statistics calculated in CLI layer where they're needed
- CLI can iterate packages to compute counts on-demand

### Why No SkippedDirs Tracking?

- `packages.Load()` doesn't provide permission error information
- Would require custom directory walking (complexity violation)
- Permission errors are rare in practice
- Acceptable trade-off for simplicity

### Why Calculate Stats in CLI?

- Parser layer stays focused: load packages
- CLI layer handles presentation logic
- Avoids coupling parser to specific statistics
- Easy to add/remove stats without changing parser
