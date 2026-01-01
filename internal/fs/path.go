// Package fs はファイルシステム操作の抽象化を提供する
package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// FilePath はファイルパスを表す型
type FilePath string

// Exists はファイルが存在するかを返す
func (f FilePath) Exists() bool {
	_, err := os.Stat(string(f))
	return !os.IsNotExist(err)
}

// Base はファイル名部分を返す
func (f FilePath) Base() string {
	return filepath.Base(string(f))
}

// DirPath はディレクトリパスを表す型
type DirPath string

// Exists はディレクトリが存在するかを返す
func (d DirPath) Exists() bool {
	info, err := os.Stat(string(d))
	return err == nil && info.IsDir()
}

// Join はパス要素を結合してFilePathを返す
func (d DirPath) Join(elem ...string) FilePath {
	return FilePath(filepath.Join(append([]string{string(d)}, elem...)...))
}

// Glob はパターンにマッチするファイルを返す
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
