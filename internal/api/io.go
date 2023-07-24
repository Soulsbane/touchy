package api

import (
	"os"
	"path"
)

type IO struct {
}

func (io *IO) CreateDirInOutputDir(name string) error {
	return os.MkdirAll(path.Join(GetOutputDir(), name), 0755)
}
