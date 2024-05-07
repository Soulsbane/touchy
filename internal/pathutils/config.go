package pathutils

import (
	"os"
	"path"
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

func GetOutputDir() string {
	dir, _ := os.Getwd()
	return dir
}

// GetAppConfigDir returns the path to the config directory for the application using companyName and applicationName.
func GetAppConfigDir() (string, error) {
	configPath, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return path.Join(configPath, companyName, applicationName), nil
}

func GetScriptsDir() string {
	configPath, err := GetAppConfigDir()

	if err != nil {
		return ""
	}

	return path.Join(configPath, "scripts")
}

func GetTemplatesDir() string {
	configPath, err := GetAppConfigDir()

	if err != nil {
		return ""
	}

	return path.Join(configPath, "templates")
}
