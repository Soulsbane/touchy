package api

import (
	"github.com/Soulsbane/touchy/internal/pathutils"
	"os"
	"path"
)

const defaultDirPermission = 0755

type IO struct {
}

func NewIO() *IO {
	return &IO{}
}

func (io *IO) CreateDirInOutputDir(name string) error {
	return os.MkdirAll(path.Join(pathutils.GetOutputDir(), name), defaultDirPermission)
}

func (io *IO) CreateDirAll(dir string) error {
	return os.MkdirAll(dir, os.FileMode(defaultDirPermission))
}
