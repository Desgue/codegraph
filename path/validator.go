package path

import (
	"fmt"
	"os"
	"path/filepath"
)

type TargetDirectory struct {
	Path string
}

func NewTargetDirectory(inputPath string) (*TargetDirectory, error) {
	var resolvedPath string

	if inputPath == "" {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		resolvedPath = currentWorkingDirectory
	} else {
		absolutePath, err := filepath.Abs(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve path '%s': %w", inputPath, err)
		}
		resolvedPath = absolutePath
	}

	canonicalPath, err := filepath.EvalSymlinks(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve symlinks for '%s': %w", resolvedPath, err)
	}

	targetDirectory := &TargetDirectory{
		Path: canonicalPath,
	}

	if err := targetDirectory.Validate(); err != nil {
		return nil, err
	}

	return targetDirectory, nil
}

func (td *TargetDirectory) Validate() error {
	fileInfo, err := os.Stat(td.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", td.Path)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied accessing '%s'", td.Path)
		}
		return fmt.Errorf("error accessing '%s': %w", td.Path, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("'%s' is a file, not a directory", td.Path)
	}

	return nil
}

func (td *TargetDirectory) String() string {
	return td.Path
}
