# Research Findings: Go Source Parser

**Date**: 2025-10-19
**Focus**: golang.org/x/tools/go/packages implementation research

## 1. Load Modes

### Decision
Use `NeedName | NeedFiles | NeedSyntax | NeedImports | NeedTypes`

### Rationale
- **NeedName**: Package name and import path for graph nodes
- **NeedFiles**: Directory location and file lists for statistics
- **NeedSyntax**: Parsed AST trees for code analysis
- **NeedImports**: Import mappings for dependency graphs
- **NeedTypes**: Required to work around NeedSyntax bug (GitHub #35331 - Syntax field isn't populated without NeedTypes)

Excluded modes:
- **NeedTypesInfo**: Expensive type-checking not needed for dependency graphs
- **NeedDeps**: Causes recursive loading of all dependencies (10x performance impact)

### Alternatives Considered
- **LoadSyntax constant**: Includes unnecessary NeedTypesInfo (10x slower)
- **LoadImports constant**: Missing AST data (NeedSyntax)
- **NeedSyntax alone**: Doesn't populate Syntax field (known bug)

## 2. Error Handling

### Decision
Partial success exits 0, only catastrophic errors exit 1

### Rationale
The `go/packages` library separates:
- **Catastrophic errors**: Returned from `packages.Load()` - pattern parsing failures, driver issues
- **Package-level errors**: Stored in `pkg.Errors` - parse errors, missing imports, syntax issues

Use `packages.PrintErrors(pkgs)` to report package errors to stderr in the format: `path/file.go:line:col: error message`

### Alternatives Considered
- **All errors exit 1**: Would prevent analysis of partially valid codebases
- **Custom error formatting**: PrintErrors already matches spec requirements

## 3. Configuration

### Decision
Minimal Config with Dir and Mode only

```go
cfg := &packages.Config{
    Mode: packages.NeedName | packages.NeedFiles |
          packages.NeedSyntax | packages.NeedImports | packages.NeedTypes,
    Dir:  targetDirectory,
}
pkgs, err := packages.Load(cfg, "./...")
```

### Rationale
- **Pattern "./..."**: Standard Go pattern for recursive package discovery
- **Dir field**: Specifies build query directory (from path.TargetDirectory)
- **No Tests**: Faster loading, test packages not needed for dependency graphs
- **No custom ParseFile**: Standard parsing is sufficient

### Alternatives Considered
- **Add Context**: Not needed for CLI tool that runs to completion
- **Tests: true**: Would add overhead, not required for current use case
- **Custom ParseFile**: Standard parser handles all requirements

## 4. Performance

### Decision
Use selective loading without NeedDeps

### Rationale
Benchmark data (GitHub #30677) shows:
- **LoadAllSyntax**: 717ms, 1.3M allocations
- **LoadImports**: 112ms, 3.4K allocations (390x reduction)
- **NeedDeps impact**: 10x performance degradation

Our configuration:
- No NeedDeps (only target directory packages)
- Minimal type-checking (enough for NeedSyntax bug)
- Pattern "./..." for efficient traversal

### Performance Target Validation
**Requirement**: 50,000 LOC in <5 seconds
- 50K LOC ≈ 100-200 files
- Estimated: 4-6 seconds (achievable)
- Can optimize if needed (caching, parallel processing)

### Alternatives Considered
- **Cache packages.Load results**: Adds complexity, not needed for single-pass CLI
- **Parallel loading**: Already handled internally by packages.Load

## 5. Module Awareness

### Decision
Rely on packages.Load automatic module detection

### Rationale
The go/packages library automatically:
- Searches for go.mod in directory tree
- Resolves import paths according to module definition
- Handles module replacements and versions
- Supports both module and GOPATH modes

No explicit configuration needed - it works transparently.

### Alternatives Considered
- **Manually parse go.mod**: Complex, unnecessary duplication
- **Use go/build**: Older API, poor module support

## 6. Return Type

### Decision
Return `[]*packages.Package` directly, no wrapper abstraction

### Rationale
Aligns with project constitution (KISS/YAGNI):
- Simplicity: No mapping layer
- Future-proof: New fields immediately available
- Testability: Use packages.Load directly in tests
- Documentation: Users reference go/packages docs

### Alternatives Considered
- **Custom Package wrapper**: Loses AST access, maintenance burden
- **Minimal interface**: Limits field access, complicates testing

## 7. Implementation Pattern

### Minimal Example

```go
package parser

import (
    "fmt"
    "golang.org/x/tools/go/packages"
)

type Result struct {
    Packages      []*packages.Package
    TotalPackages int
    TotalFiles    int
    ErrorCount    int
    SkippedDirs   []string
}

func Load(targetDir string) (*Result, error) {
    cfg := &packages.Config{
        Mode: packages.NeedName | packages.NeedFiles |
              packages.NeedSyntax | packages.NeedImports | packages.NeedTypes,
        Dir: targetDir,
    }

    pkgs, err := packages.Load(cfg, "./...")
    if err != nil {
        return nil, fmt.Errorf("failed to load packages: %w", err)
    }

    errorCount := packages.PrintErrors(pkgs)

    result := &Result{
        Packages:      pkgs,
        TotalPackages: len(pkgs),
        ErrorCount:    errorCount,
    }

    for _, pkg := range pkgs {
        result.TotalFiles += len(pkg.GoFiles)
    }

    return result, nil
}
```

## 8. Known Limitations

1. **NeedSyntax Bug**: Must include NeedTypes even though type-checking isn't needed (GitHub #35331)
   - Mitigation: Accept performance cost, unavoidable

2. **Performance Target**: 5 seconds for 50K LOC is tight
   - Mitigation: Can add caching if needed

3. **LoadMode Bugs**: Documentation mentions open bugs with mode bit interactions
   - Mitigation: Use well-tested combinations

## 9. References

- **Documentation**: https://pkg.go.dev/golang.org/x/tools/go/packages
- **Source**: https://github.com/golang/tools/tree/master/go/packages
- **Key Issues**: #30677 (performance), #35331 (NeedSyntax bug)
- **Examples**: staticcheck, gopls, structlayout
