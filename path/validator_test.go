package path

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTargetDirectory_ValidAbsolutePath(t *testing.T) {
	tempDir := t.TempDir()

	targetDirectory, err := NewTargetDirectory(tempDir)

	require.NoError(t, err)
	expectedPath, _ := filepath.EvalSymlinks(tempDir)
	actualPath, _ := filepath.EvalSymlinks(targetDirectory.Path)
	assert.Equal(t, expectedPath, actualPath)
}

func TestNewTargetDirectory_ValidRelativePath(t *testing.T) {
	tempDir := t.TempDir()
	originalWorkingDirectory, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWorkingDirectory)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	subDir := "testsubdir"
	err = os.Mkdir(filepath.Join(tempDir, subDir), 0755)
	require.NoError(t, err)

	targetDirectory, err := NewTargetDirectory(subDir)

	require.NoError(t, err)
	expectedPath, _ := filepath.EvalSymlinks(filepath.Join(tempDir, subDir))
	actualPath, _ := filepath.EvalSymlinks(targetDirectory.Path)
	assert.Equal(t, expectedPath, actualPath)
}

func TestNewTargetDirectory_CurrentDirectoryDefault(t *testing.T) {
	currentWorkingDirectory, err := os.Getwd()
	require.NoError(t, err)

	targetDirectory, err := NewTargetDirectory("")

	require.NoError(t, err)
	assert.Equal(t, currentWorkingDirectory, targetDirectory.Path)
}

func TestNewTargetDirectory_NonexistentDirectory(t *testing.T) {
	nonexistentPath := "/nonexistent/path/that/should/not/exist"

	targetDirectory, err := NewTargetDirectory(nonexistentPath)

	assert.Nil(t, targetDirectory)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "directory does not exist")
	assert.Contains(t, err.Error(), nonexistentPath)
}

func TestNewTargetDirectory_FileInsteadOfDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err := os.WriteFile(tempFile, []byte("test content"), 0644)
	require.NoError(t, err)

	targetDirectory, err := NewTargetDirectory(tempFile)

	assert.Nil(t, targetDirectory)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is a file, not a directory")
	assert.Contains(t, err.Error(), tempFile)
}

func TestNewTargetDirectory_SymbolicLinkToDirectory(t *testing.T) {
	tempDir := t.TempDir()
	realDir := filepath.Join(tempDir, "realdir")
	err := os.Mkdir(realDir, 0755)
	require.NoError(t, err)

	symlinkPath := filepath.Join(tempDir, "linkdir")
	err = os.Symlink(realDir, symlinkPath)
	require.NoError(t, err)

	targetDirectory, err := NewTargetDirectory(symlinkPath)

	require.NoError(t, err)
	expectedPath, _ := filepath.EvalSymlinks(symlinkPath)
	actualPath, _ := filepath.EvalSymlinks(targetDirectory.Path)
	assert.Equal(t, expectedPath, actualPath)
}

func TestNewTargetDirectory_SymbolicLinkToFile(t *testing.T) {
	tempDir := t.TempDir()
	realFile := filepath.Join(tempDir, "realfile.txt")
	err := os.WriteFile(realFile, []byte("test content"), 0644)
	require.NoError(t, err)

	symlinkPath := filepath.Join(tempDir, "linkfile")
	err = os.Symlink(realFile, symlinkPath)
	require.NoError(t, err)

	targetDirectory, err := NewTargetDirectory(symlinkPath)

	assert.Nil(t, targetDirectory)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is a file, not a directory")
}

func TestNewTargetDirectory_PathWithSpaces(t *testing.T) {
	tempDir := t.TempDir()
	dirWithSpaces := filepath.Join(tempDir, "dir with spaces")
	err := os.Mkdir(dirWithSpaces, 0755)
	require.NoError(t, err)

	targetDirectory, err := NewTargetDirectory(dirWithSpaces)

	require.NoError(t, err)
	expectedPath, _ := filepath.EvalSymlinks(dirWithSpaces)
	actualPath, _ := filepath.EvalSymlinks(targetDirectory.Path)
	assert.Equal(t, expectedPath, actualPath)
}

func TestTargetDirectory_String(t *testing.T) {
	tempDir := t.TempDir()

	targetDirectory, err := NewTargetDirectory(tempDir)
	require.NoError(t, err)

	result := targetDirectory.String()

	expectedPath, _ := filepath.EvalSymlinks(tempDir)
	actualPath, _ := filepath.EvalSymlinks(result)
	assert.Equal(t, expectedPath, actualPath)
}
