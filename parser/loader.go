package parser

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

// Load parses all Go packages in targetDir and returns them.
// Returns error only for catastrophic failures (pattern parsing, driver issues).
// Package-level parse errors are printed to stderr via packages.PrintErrors().
// Each returned package contains Errors field with parse failures.
//
// Comment Access Patterns:
// - Package-level comments: pkg.Syntax[i].Doc
// - All comment nodes in file: pkg.Syntax[i].Comments
// - Function/type comments: Access via ast.Walk on pkg.Syntax[i]
// Comments are preserved with NeedSyntax flag for future documentation analysis.
func Load(targetDir string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedSyntax | packages.NeedImports | packages.NeedTypes,
		Dir:   targetDir,
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	packages.PrintErrors(pkgs)

	return pkgs, nil
}
