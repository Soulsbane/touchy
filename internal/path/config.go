package path

import (
	"os"
	"path/filepath"
)

const companyName = "Raijinsoft"
const applicationName = "touchy"

// GetAppConfigDir returns the path to the config directory for the application using companyName and applicationName.
func GetAppConfigDir() (string, error) {
	path, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(path, companyName, applicationName), nil
}
