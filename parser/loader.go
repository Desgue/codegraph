package parser

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Load parses all Go packages in targetDir and returns them with error count.
// Returns error only for catastrophic failures (pattern parsing, driver issues).
// Package-level parse errors are printed to stderr via packages.PrintErrors().
// Each returned package contains Errors field with parse failures.
// The error count returned is the number of packages with errors (from packages.PrintErrors).
//
// Deduplication: When includeTests is true, go/packages returns both regular and test
// variants of each package. This function deduplicates by keeping only the variant
// with the most files (which includes both production and test files).
// Synthetic .test packages are filtered out.
//
// Comment Access Patterns:
// - Package-level comments: pkg.Syntax[i].Doc
// - All comment nodes in file: pkg.Syntax[i].Comments
// - Function/type comments: Access via ast.Walk on pkg.Syntax[i]
// Comments are preserved with NeedSyntax flag for future documentation analysis.
func Load(targetDir string, includeTests bool) ([]*packages.Package, int, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedSyntax | packages.NeedImports | packages.NeedTypes,
		Dir:   targetDir,
		Tests: includeTests,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load packages: %w", err)
	}

	errorCount := packages.PrintErrors(pkgs)

	// Deduplicate packages and filter synthetic test packages
	deduplicated := deduplicatePackages(pkgs)

	// Sort packages by import path for deterministic output
	sort.Slice(deduplicated, func(i, j int) bool {
		return deduplicated[i].PkgPath < deduplicated[j].PkgPath
	})

	return deduplicated, errorCount, nil
}

// deduplicatePackages removes duplicate package variants and synthetic test packages.
// When Tests is true, go/packages returns multiple variants of the same package.
// This keeps only the variant with the most files (test variant has production + test files).
func deduplicatePackages(pkgs []*packages.Package) []*packages.Package {
	seen := make(map[string]*packages.Package)

	for _, pkg := range pkgs {
		// Skip synthetic test binary packages (e.g., package.test)
		if strings.HasSuffix(pkg.PkgPath, ".test") {
			continue
		}

		if existing, found := seen[pkg.PkgPath]; found {
			// Keep the variant with more files (test variant includes production + test)
			if len(pkg.GoFiles) > len(existing.GoFiles) {
				seen[pkg.PkgPath] = pkg
			}
		} else {
			seen[pkg.PkgPath] = pkg
		}
	}

	// Convert map back to slice
	result := make([]*packages.Package, 0, len(seen))
	for _, pkg := range seen {
		result = append(result, pkg)
	}

	return result
}
