package safety

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Return whether the given path (dir/file/etc) currently exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Return whether the given path exists and is a directory.
func IsDir(path string) bool {
	if !PathExists(path) {
		return false
	}
	stat, err := os.Stat(path)
	if err != nil {
		panic(err.Error()) // Was just successful so should be impossible.
	}
	return stat.IsDir()
}

// Return whether the given path exists and is a regular file.
func IsFile(path string) bool {
	if !PathExists(path) {
		return false
	}
	return FileMode(path).IsRegular()
}

// Simple error-free wrapper to get the file-mode data.
//
// WARNING: It is the caller's responsibility to ensure that the path exists in
// the first-place and if shirked, may panic.
func FileMode(path string) os.FileMode {
	if !PathExists(path) {
		panic("BAD CODE-PATH ALLOWED fileio.FileMode() call without fileio.IsXXX() or fileio.PathExists() CALL")
	}

	stat, err := os.Stat(path)
	if err != nil {
		panic(err.Error()) // Was just successful so should be impossible.
	}

	return stat.Mode()
}

// Get a list of all the files under the given root directory.
func GetFileListing(rootDir string) ([]string, error) {
	files := []string{}
	walkFunc := func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		switch {
		case entry.IsDir():
			return nil
		case FileMode(path).IsRegular():
			files = append(files, path)
		default:
			return fmt.Errorf("unhandled filetype. file: %q", path)
		}

		return nil
	}

	if err := filepath.WalkDir(rootDir, walkFunc); err != nil {
		return nil, fmt.Errorf("filesystem walk failed: %v", err)
	}

	return files, nil
}

// Simple IO-wrapper to copy a file from one place to another.
func CopyFile(fromPath string, toPath string) error {
	toDir, _ := filepath.Split(toPath)
	switch {
	case !PathExists(fromPath):
		return errors.New("from-path doesn't exist")
	case !IsFile(fromPath):
		return errors.New("from-path isn't a standard file")
	case !PathExists(toDir):
		return errors.New("to-path's parent directory doesn't exist")
	case PathExists(toPath):
		return errors.New("to-path already exist")
	}

	data, err := os.ReadFile(fromPath)
	if err != nil {
		return fmt.Errorf("could not read from-path: %w", err)
	}

	if err := os.WriteFile(toPath, data, FileMode(fromPath)); err != nil {
		return fmt.Errorf("could not write to-path: %w", err)
	}

	return nil
}
