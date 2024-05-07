package api

import (
	"github.com/Soulsbane/touchy/internal/pathutils"
	"os"
	"path"
)

type IO struct {
}

func (io *IO) CreateDirInOutputDir(name string) error {
	return os.MkdirAll(path.Join(pathutils.GetOutputDir(), name), 0755)
}
