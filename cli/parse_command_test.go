package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewParseCommand(t *testing.T) {
	validTests := []struct {
		name             string
		setup            func(t *testing.T) []string
		wantOutputFile   string
		wantIncludeTests bool
	}{
		{
			name: "valid minimal args with current directory",
			setup: func(t *testing.T) []string {
				return []string{"--output", "out.graphml"}
			},
			wantOutputFile:   "out.graphml",
			wantIncludeTests: false,
		},
		{
			name: "valid args with directory",
			setup: func(t *testing.T) []string {
				return []string{"--output", "out.graphml", t.TempDir()}
			},
			wantOutputFile:   "out.graphml",
			wantIncludeTests: false,
		},
		{
			name: "valid args with include tests flag",
			setup: func(t *testing.T) []string {
				return []string{"--output", "out.graphml", "--include-tests"}
			},
			wantOutputFile:   "out.graphml",
			wantIncludeTests: true,
		},
		{
			name: "valid args with explicit include tests false",
			setup: func(t *testing.T) []string {
				return []string{"--output", "out.graphml", "--include-tests=false"}
			},
			wantOutputFile:   "out.graphml",
			wantIncludeTests: false,
		},
		{
			name: "valid args with relative directory",
			setup: func(t *testing.T) []string {
				tempDir := t.TempDir()
				originalWd, _ := os.Getwd()
				t.Cleanup(func() { os.Chdir(originalWd) })

				if err := os.Chdir(tempDir); err != nil {
					t.Fatalf("failed to change directory: %v", err)
				}

				subDir := filepath.Join(tempDir, "subdir")
				if err := os.Mkdir(subDir, 0755); err != nil {
					t.Fatalf("failed to create subdirectory: %v", err)
				}

				return []string{"--output", "out.graphml", "./subdir"}
			},
			wantOutputFile:   "out.graphml",
			wantIncludeTests: false,
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.setup(t)
			cmd, err := NewParseCommand(args)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if cmd.OutputFile != tt.wantOutputFile {
				t.Errorf("OutputFile = %q, want %q", cmd.OutputFile, tt.wantOutputFile)
			}
			if cmd.IncludeTests != tt.wantIncludeTests {
				t.Errorf("IncludeTests = %v, want %v", cmd.IncludeTests, tt.wantIncludeTests)
			}
			if cmd.TargetDirectory == nil {
				t.Fatal("expected TargetDirectory to be set")
			}
		})
	}

	tests := []struct {
		name string
		args []string
		setup func(t *testing.T) []string
	}{
		{
			name: "missing output flag returns error",
			args: []string{},
		},
		{
			name: "empty output flag returns error",
			args: []string{"--output", ""},
		},
		{
			name: "non-existent directory returns error",
			args: []string{"--output", "out.graphml", "/non/existent/path"},
		},
		{
			name: "file instead of directory returns error",
			setup: func(t *testing.T) []string {
				tempDir := t.TempDir()
				tempFile := filepath.Join(tempDir, "file.txt")
				if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				return []string{"--output", "out.graphml", tempFile}
			},
		},
		{
			name: "unknown flag returns error",
			args: []string{"--output", "out.graphml", "--unknown-flag"},
		},
		{
			name: "invalid boolean syntax returns error",
			args: []string{"--output", "out.graphml", "--include-tests=invalid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args
			if tt.setup != nil {
				args = tt.setup(t)
			}

			_, err := NewParseCommand(args)
			if err == nil {
				t.Fatal("expected error but got none")
			}
		})
	}
}

func TestParseCommand_Validate(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) *ParseCommand
		wantError bool
	}{
		{
			name: "valid command passes validation",
			setup: func(t *testing.T) *ParseCommand {
				cmd, err := NewParseCommand([]string{"--output", "out.graphml", t.TempDir()})
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return cmd
			},
			wantError: false,
		},
		{
			name: "missing output file fails validation",
			setup: func(t *testing.T) *ParseCommand {
				validCmd, err := NewParseCommand([]string{"--output", "temp.graphml", t.TempDir()})
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				return &ParseCommand{
					TargetDirectory: validCmd.TargetDirectory,
					OutputFile:      "",
					IncludeTests:    false,
				}
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup(t)
			err := cmd.Validate()

			if tt.wantError && err == nil {
				t.Fatal("expected validation error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no validation error, got %v", err)
			}
		})
	}
}

func TestParseCommand_Execute(t *testing.T) {
	t.Run("returns no error", func(t *testing.T) {
		cmd, err := NewParseCommand([]string{"--output", "out.graphml"})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		if err := cmd.Execute(); err != nil {
			t.Errorf("expected no error from Execute, got %v", err)
		}
	})
}
