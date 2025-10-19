# Quickstart: Go Source Parser Implementation

**Date**: 2025-10-19
**Feature**: 002-go-source-parser
**Estimated Time**: 2-3 hours

## Prerequisites

- Go 1.24+ installed
- Existing codegraph project structure
- Feature branch `002-go-source-parser` checked out

## Step 1: Add Dependency (5 minutes)

```bash
go get golang.org/x/tools/go/packages
```

**Verify**:
```bash
go mod tidy
cat go.mod  # Should include golang.org/x/tools
```

## Step 2: Create Parser Package (10 minutes)

### Create package directory
```bash
mkdir parser
```

### Create loader.go

```go
// parser/loader.go
package parser

import (
    "fmt"
    "golang.org/x/tools/go/packages"
)

// Load parses all Go packages in targetDir and returns them.
// Returns error only for catastrophic failures.
// Package-level parse errors are printed to stderr via packages.PrintErrors().
func Load(targetDir string) ([]*packages.Package, error) {
    cfg := &packages.Config{
        Mode: packages.NeedName | packages.NeedFiles |
              packages.NeedSyntax | packages.NeedImports,
        Dir: targetDir,
    }

    pkgs, err := packages.Load(cfg, "./...")
    if err != nil {
        return nil, fmt.Errorf("failed to load packages: %w", err)
    }

    packages.PrintErrors(pkgs)

    return pkgs, nil
}
```

**Test**:
```bash
go build ./parser
```

## Step 3: Write Tests (30 minutes)

### Create test fixtures

```bash
mkdir -p parser/testdata/valid-project
mkdir -p parser/testdata/empty-dir
mkdir -p parser/testdata/syntax-errors
```

### Create valid test project

```go
// parser/testdata/valid-project/main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Create syntax error test case

```go
// parser/testdata/syntax-errors/invalid.go
package main

func main() {
    fmt.Println("missing import and semicolon"
}
```

### Create loader_test.go

```go
// parser/loader_test.go
package parser

import (
    "path/filepath"
    "testing"
)

func TestLoad_ValidProject(t *testing.T) {
    dir, _ := filepath.Abs("testdata/valid-project")

    pkgs, err := Load(dir)

    if err != nil {
        t.Fatalf("Load() error = %v", err)
    }
    if len(pkgs) != 1 {
        t.Errorf("len(pkgs) = %d, want 1", len(pkgs))
    }
    if len(pkgs[0].GoFiles) != 1 {
        t.Errorf("len(pkgs[0].GoFiles) = %d, want 1", len(pkgs[0].GoFiles))
    }
}

func TestLoad_EmptyDirectory(t *testing.T) {
    dir, _ := filepath.Abs("testdata/empty-dir")

    pkgs, err := Load(dir)

    if err != nil {
        t.Fatalf("Load() error = %v", err)
    }
    if len(pkgs) != 0 {
        t.Errorf("len(pkgs) = %d, want 0", len(pkgs))
    }
}

func TestLoad_WithSyntaxErrors(t *testing.T) {
    dir, _ := filepath.Abs("testdata/syntax-errors")

    pkgs, err := Load(dir)

    // Should NOT return error for parse failures
    if err != nil {
        t.Fatalf("Load() error = %v", err)
    }
    // Should have packages with errors
    if len(pkgs) == 0 {
        t.Fatal("len(pkgs) = 0, want > 0")
    }
    if len(pkgs[0].Errors) == 0 {
        t.Error("len(pkgs[0].Errors) = 0, want > 0")
    }
}
```

**Test**:
```bash
go test ./parser -v
```

## Step 4: Integrate with CLI (20 minutes)

### Read existing ParseCommand

```bash
cat cli/parse.go
```

### Update cli/parse.go Execute method

```go
// cli/parse.go
func (c *ParseCommand) Execute() error {
    targetDir, err := path.NewTargetDirectory(c.targetPath)
    if err != nil {
        return err
    }

    // NEW: Parse packages
    pkgs, err := parser.Load(targetDir.Path())
    if err != nil {
        return err
    }

    // NEW: Calculate statistics
    totalPackages := len(pkgs)
    totalFiles := 0
    errorCount := 0
    for _, pkg := range pkgs {
        totalFiles += len(pkg.GoFiles)
        errorCount += len(pkg.Errors)
    }

    // NEW: Print statistics
    fmt.Printf("Loaded %d packages, parsed %d files\n", totalPackages, totalFiles)
    if errorCount > 0 {
        fmt.Fprintf(os.Stderr, "Encountered %d parse errors\n", errorCount)
    }

    return nil
}
```

**Test**:
```bash
make build
./bin/codegraph parse .
```

Expected output:
```
Loaded 4 packages, parsed 8 files
```

## Step 5: Write CLI Tests (30 minutes)

### Create cli/parse_test.go tests

```go
// cli/parse_test.go
func TestParseCommand_Execute_ValidProject(t *testing.T) {
    cmd := &ParseCommand{targetPath: "../parser/testdata/valid-project"}

    err := cmd.Execute()

    if err != nil {
        t.Fatalf("Execute() error = %v", err)
    }
}

func TestParseCommand_Execute_EmptyDirectory(t *testing.T) {
    cmd := &ParseCommand{targetPath: "../parser/testdata/empty-dir"}

    err := cmd.Execute()

    if err != nil {
        t.Fatalf("Execute() error = %v", err)
    }
}
```

**Test**:
```bash
go test ./cli -v
```

## Step 6: Manual Testing (15 minutes)

### Test on real project

```bash
make build

# Test on self
./bin/codegraph parse .

# Test on empty directory
mkdir /tmp/empty
./bin/codegraph parse /tmp/empty

# Test on project with errors
mkdir -p /tmp/badgo
echo "package main\nfunc main() { missing brace" > /tmp/badgo/bad.go
./bin/codegraph parse /tmp/badgo
```

### Verify output format

Expected for valid project:
```
Loaded N packages, parsed M files
```

Expected for project with errors (stderr first, then stdout):
```
path/to/file.go:line:col: error message
Loaded N packages, parsed M files
Encountered X parse errors
```

## Step 7: Update Documentation (10 minutes)

### Update CLAUDE.md

Add parser package description:
```markdown
### Core Structure

- **parser/**: Go source code parsing
  - `Load()`: Parse packages using go/packages, returns []*packages.Package
```

**Test**:
```bash
git diff CLAUDE.md
```

## Step 8: Run Full Test Suite (10 minutes)

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Build binary
make build

# Test binary
./bin/codegraph parse .
```

## Verification Checklist

- [ ] Dependency added: `golang.org/x/tools/go/packages`
- [ ] Parser package created: `parser/loader.go`
- [ ] Tests pass: `go test ./parser`
- [ ] CLI integrated: `cli/parse.go` uses `parser.Load()`
- [ ] CLI calculates statistics from returned packages
- [ ] CLI tests pass: `go test ./cli`
- [ ] Manual testing successful on valid, empty, and error cases
- [ ] Documentation updated: `CLAUDE.md`
- [ ] Build succeeds: `make build`
- [ ] Binary works: `./bin/codegraph parse .`

## Common Issues

### Issue: "package golang.org/x/tools/go/packages not found"
**Solution**: Run `go get golang.org/x/tools/go/packages` and `go mod tidy`

### Issue: Tests fail with "no such file or directory"
**Solution**: Ensure testdata directories exist and have correct structure

### Issue: "Syntax field is nil"
**Solution**: Verify LoadMode includes `NeedSyntax` and does NOT include `NeedTypes`

### Issue: Performance is slow
**Solution**: Verify LoadMode does NOT include `NeedDeps` or `NeedTypes` (causes 10x slowdown)

## Next Steps

After completing this implementation:

1. Run `/speckit.tasks` to generate task breakdown
2. Commit work: `git add . && git commit -m "feat(parser): add Go source code parsing"`
3. Continue to graph construction phase (future feature)

## Estimated Completion Time

- Experienced Go developer: 2 hours
- Intermediate Go developer: 3 hours
- New to go/packages: 4 hours

## Success Criteria

- [ ] `go test ./...` passes
- [ ] `make build` succeeds
- [ ] `./bin/codegraph parse .` shows "Loaded N packages, parsed M files"
- [ ] Parse errors are reported to stderr with file:line:col format
- [ ] Empty directories show "Loaded 0 packages, parsed 0 files" and exit 0
- [ ] Projects with syntax errors show error count but exit 0
