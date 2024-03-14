package path

import (
	"os"
	"path/filepath"
)

// CleanPath taken from https://github.com/opencontainers/runc/blob/main/libcontainer/utils/utils.go
func CleanPath(path string) string {
	if path == "" {
		return ""
	}

	path = filepath.Clean(path)

	if !filepath.IsAbs(path) {
		path = filepath.Clean(string(os.PathSeparator) + path)
		path, _ = filepath.Rel(string(os.PathSeparator), path)
	}

	return filepath.Clean(path)
}
