package api

import (
	"errors"
	"fmt"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"os"
	"path"
)

const defaultDirPermission = 0755

var ErrFailedToCreateOutputDir = errors.New("failed to create directory in output directory")
var ErrFailedToCreateDirPath = errors.New("failed to create directory path")

type IO struct {
}

func NewIO() *IO {
	return &IO{}
}

func (io *IO) CreateDirInOutputDir(name string) error {
	outputPath := path.Join(pathutils.GetOutputDir(), name)
	err := os.MkdirAll(outputPath, defaultDirPermission)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateOutputDir, err)
	}

	return nil
}

func (io *IO) CreateDirAll(dir string) error {
	err := os.MkdirAll(dir, defaultDirPermission)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateDirPath, err)
	}

	return nil
}
