package path

import (
	"os"
	"path/filepath"
)

const companyName = "Raijinsoft"
const applicationName = "touchy"

func SetupConfigDir() {
	configScriptsDir := GetScriptsDir()

	if _, err := os.Stat(configScriptsDir); os.IsNotExist(err) || configScriptsDir == "" {
		os.MkdirAll(configScriptsDir, os.ModePerm)
	}
}

// GetAppConfigDir returns the path to the config directory for the application using companyName and applicationName.
func GetAppConfigDir() (string, error) {
	path, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(path, companyName, applicationName), nil
}

func GetScriptsDir() string {
	path, err := GetAppConfigDir()

	if err != nil {
		return ""
	}

	return filepath.Join(path, "scripts")
}
