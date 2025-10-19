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

	pkgs, errorCount, err := Load(testDir, true)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
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

func TestLoad_EmptyDirectory(t *testing.T) {
	testDir := t.TempDir()

	goMod := filepath.Join(testDir, "go.mod")
	modContent := "module testmod\n\ngo 1.24\n"
	if err := os.WriteFile(goMod, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	pkgs, errorCount, err := Load(testDir, true)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}

	if len(pkgs) != 0 {
		t.Errorf("Expected 0 packages in empty directory, got %d", len(pkgs))
	}
}

func TestLoad_MultiplePackages(t *testing.T) {
	testDir := t.TempDir()

	goMod := filepath.Join(testDir, "go.mod")
	modContent := "module testmod\n\ngo 1.24\n"
	if err := os.WriteFile(goMod, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create package 1
	pkg1Dir := filepath.Join(testDir, "pkg1")
	if err := os.Mkdir(pkg1Dir, 0755); err != nil {
		t.Fatalf("Failed to create pkg1 directory: %v", err)
	}
	pkg1File := filepath.Join(pkg1Dir, "main.go")
	pkg1Content := "package pkg1\n\nconst Name = \"pkg1\"\n"
	if err := os.WriteFile(pkg1File, []byte(pkg1Content), 0644); err != nil {
		t.Fatalf("Failed to create pkg1/main.go: %v", err)
	}

	// Create package 2
	pkg2Dir := filepath.Join(testDir, "pkg2")
	if err := os.Mkdir(pkg2Dir, 0755); err != nil {
		t.Fatalf("Failed to create pkg2 directory: %v", err)
	}
	pkg2File := filepath.Join(pkg2Dir, "main.go")
	pkg2Content := "package pkg2\n\nconst Name = \"pkg2\"\n"
	if err := os.WriteFile(pkg2File, []byte(pkg2Content), 0644); err != nil {
		t.Fatalf("Failed to create pkg2/main.go: %v", err)
	}

	pkgs, errorCount, err := Load(testDir, true)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}

	if len(pkgs) != 2 {
		t.Fatalf("Expected 2 packages, got %d", len(pkgs))
	}

	// Verify packages are sorted by import path
	if pkgs[0].PkgPath > pkgs[1].PkgPath {
		t.Errorf("Packages not sorted: %s should come before %s", pkgs[0].PkgPath, pkgs[1].PkgPath)
	}
}

func TestLoad_WithSyntaxErrors(t *testing.T) {
	testDir := t.TempDir()

	goMod := filepath.Join(testDir, "go.mod")
	modContent := "module testmod\n\ngo 1.24\n"
	if err := os.WriteFile(goMod, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	invalidFile := filepath.Join(testDir, "invalid.go")
	invalidContent := "package main\n\nfunc broken( {\n"
	if err := os.WriteFile(invalidFile, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create invalid.go: %v", err)
	}

	pkgs, errorCount, err := Load(testDir, true)
	// Should not return error for parse errors (partial failure)
	if err != nil {
		t.Fatalf("Load() should not error on syntax errors, got %v", err)
	}

	// Should report errors via errorCount
	if errorCount == 0 {
		t.Error("Expected errorCount > 0 for syntax errors")
	}

	// Should still return the package despite errors
	if len(pkgs) != 1 {
		t.Errorf("Expected 1 package with errors, got %d", len(pkgs))
	}
}

func TestLoad_DeduplicationWithTests(t *testing.T) {
	testDir := t.TempDir()

	goMod := filepath.Join(testDir, "go.mod")
	modContent := "module testmod\n\ngo 1.24\n"
	if err := os.WriteFile(goMod, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create production file
	prodFile := filepath.Join(testDir, "main.go")
	prodContent := "package main\n\nfunc Hello() string { return \"hello\" }\n"
	if err := os.WriteFile(prodFile, []byte(prodContent), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Create test file
	testFile := filepath.Join(testDir, "main_test.go")
	testContent := "package main\n\nimport \"testing\"\n\nfunc TestHello(t *testing.T) {}\n"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create main_test.go: %v", err)
	}

	pkgs, errorCount, err := Load(testDir, true)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}

	// Should deduplicate to single package containing both production and test files
	if len(pkgs) != 1 {
		t.Fatalf("Expected 1 deduplicated package, got %d", len(pkgs))
	}

	pkg := pkgs[0]
	// Test variant should have both files
	if len(pkg.GoFiles) < 2 {
		t.Errorf("Expected at least 2 files (production + test), got %d", len(pkg.GoFiles))
	}
}

func TestLoad_FiltersSyntheticTestPackages(t *testing.T) {
	testDir := t.TempDir()

	goMod := filepath.Join(testDir, "go.mod")
	modContent := "module testmod\n\ngo 1.24\n"
	if err := os.WriteFile(goMod, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	mainFile := filepath.Join(testDir, "main.go")
	mainContent := "package main\n\nfunc main() {}\n"
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	testFile := filepath.Join(testDir, "main_test.go")
	testContent := "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}\n"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create main_test.go: %v", err)
	}

	pkgs, errorCount, err := Load(testDir, true)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}

	// Verify no .test packages in results
	for _, pkg := range pkgs {
		if pkg.PkgPath == "testmod.test" || pkg.Name == "main.test" {
			t.Errorf("Synthetic .test package should be filtered out: %s", pkg.PkgPath)
		}
	}
}

func TestLoad_IncludeTestsFalse(t *testing.T) {
	testDir := t.TempDir()

	goMod := filepath.Join(testDir, "go.mod")
	modContent := "module testmod\n\ngo 1.24\n"
	if err := os.WriteFile(goMod, []byte(modContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	mainFile := filepath.Join(testDir, "main.go")
	mainContent := "package main\n\nfunc main() {}\n"
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	testFile := filepath.Join(testDir, "main_test.go")
	testContent := "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}\n"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create main_test.go: %v", err)
	}

	pkgs, errorCount, err := Load(testDir, false)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}

	if len(pkgs) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(pkgs))
	}

	pkg := pkgs[0]
	// When includeTests is false, should only have production file
	hasTestFile := false
	for _, file := range pkg.GoFiles {
		if filepath.Base(file) == "main_test.go" {
			hasTestFile = true
			break
		}
	}

	if hasTestFile {
		t.Error("Test files should not be included when includeTests is false")
	}
}
