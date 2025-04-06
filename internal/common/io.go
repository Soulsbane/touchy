package common

import (
	"fmt"
	"os"
)

func GetFileData(path string) ([]byte, error) {
	if data, err := os.ReadFile(path); err != nil {
		return data, fmt.Errorf("%w: %w", ErrFailedToReadFile, err)
	} else {
		return data, nil
	}
}
