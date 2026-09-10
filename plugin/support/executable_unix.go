//go:build !windows

package support

import "os"

func IsExecutable(filePath string, fileInfo os.FileInfo) bool {
	return fileInfo.Mode().IsRegular() && fileInfo.Mode()&0111 != 0
}
