package templates

import (
	"errors"
	"log"

	"github.com/pelletier/go-toml/v2"
)

type Language struct {
	DirName         string                    // Name of the directory under the templates directory.
	Info            LanguageInfo              // Each language has a config file in its root directory call config.toml
	TemplateConfigs map[string]TemplateConfig // A list of all the templates in the language directory.
}

type BaseConfig struct {
	Name                  string
	DefaultOutputFileName string
	Description           string
	Extension             string
}

type TemplateConfig struct {
	BaseConfig
}

type LanguageInfo struct {
	BaseConfig
}

func loadLanguageInfoFile(languageName string) LanguageInfo {
	data, err := templatesDir.ReadFile(languageName)
	config := LanguageInfo{}

	if err != nil {
		log.Fatal(errors.New("Failed to load config file: " + languageName))
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(errors.New("Failed to read config file: " + languageName))
	}

	return config
}

func loadLanguageConfigFile(languageName string) TemplateConfig {
	data, err := templatesDir.ReadFile(languageName)
	config := TemplateConfig{}

	if err != nil {
		log.Fatal(errors.New("Failed to load config file: " + languageName))
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(errors.New("Failed to read config file: " + languageName))
	}

	return config
}
