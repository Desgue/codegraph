package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_PreservesComments(t *testing.T) {
	testDir := t.TempDir()

	goMod := filepath.Join(testDir, "go.mod")
	modContent := "module testmod\n\ngo 1.24\n"
	if err := os.WriteFile(goMod, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	testFile := filepath.Join(testDir, "main.go")
	content := `// Package main provides the entry point.
package main

// Greeting is a constant message.
const Greeting = "Hello, World!"

// main is the entry point function.
func main() {
	println(Greeting)
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	pkgs, err := Load(testDir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(pkgs))
	}

	pkg := pkgs[0]
	if len(pkg.Syntax) != 1 {
		t.Fatalf("Expected 1 AST file, got %d", len(pkg.Syntax))
	}

	astFile := pkg.Syntax[0]

	// Verify package-level comment
	if astFile.Doc == nil {
		t.Error("Package-level comment not preserved")
	} else if len(astFile.Doc.List) == 0 {
		t.Error("Package-level comment list is empty")
	}

	// Verify comment nodes exist in AST
	if len(astFile.Comments) == 0 {
		t.Error("No comment nodes found in AST")
	}
}
