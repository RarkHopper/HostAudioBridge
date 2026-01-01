package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

type FilePath string

func (f FilePath) Exists() bool {
	_, err := os.Stat(string(f))
	return !os.IsNotExist(err)
}

func (f FilePath) Base() string {
	return filepath.Base(string(f))
}

type DirPath string

func (d DirPath) Exists() bool {
	info, err := os.Stat(string(d))
	return err == nil && info.IsDir()
}

func (d DirPath) Join(elem ...string) FilePath {
	return FilePath(filepath.Join(append([]string{string(d)}, elem...)...))
}

func (d DirPath) Glob(pattern string) ([]FilePath, error) {
	matches, err := filepath.Glob(filepath.Join(string(d), pattern))
	if err != nil {
		return nil, fmt.Errorf("glob failed: %w", err)
	}
	result := make([]FilePath, len(matches))
	for i, m := range matches {
		result[i] = FilePath(m)
	}
	return result, nil
}
