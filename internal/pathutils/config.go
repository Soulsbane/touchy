package pathutils

import (
	"os"
	"path/filepath"
)

const companyName = "Raijinsoft"
const applicationName = "touchy"

func SetupConfigDir() error {
	configScriptsDir := GetScriptsDir()

	if _, err := os.Stat(configScriptsDir); os.IsNotExist(err) || configScriptsDir == "" {
		return os.MkdirAll(configScriptsDir, os.ModePerm)
	}

	return nil
}

// GetAppConfigDir returns the pathutils to the config directory for the application using companyName and applicationName.
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

func GetTemplatesDir() string {
	path, err := GetAppConfigDir()

	if err != nil {
		return ""
	}

	return filepath.Join(path, "templates")
}
