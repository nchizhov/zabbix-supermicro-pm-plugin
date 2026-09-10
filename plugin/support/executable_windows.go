//go:build windows

package support

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func IsExecutable(filePath string, fileInfo os.FileInfo) bool {
	if !fileInfo.Mode().IsRegular() {
		return false
	}
	ext := strings.ToUpper(filepath.Ext(filePath))
	if ext == "" {
		return false
	}

	pathext := os.Getenv("PATHEXT")
	if pathext == "" {
		pathext = ".COM;.EXE;.BAT;.CMD"
	}
	pathexts := strings.Split(strings.ToUpper(pathext), ";")
	return slices.Contains(pathexts, ext)
}
