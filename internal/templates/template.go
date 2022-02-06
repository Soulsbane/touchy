package templates

import (
	"errors"
	"log"

	"github.com/pelletier/go-toml/v2"
)

type Language struct {
	DirName       string         // Name of the directory under the templates directory.
	Config        LanguageConfig // Each language has a config file in its root directory call config.toml
	TemplatesInfo []TemplateInfo // A list of all the templates in the language directory.
}

type LanguageConfig struct {
	Name                  string
	DefaultOutputFileName string
	Description           string
	Extension             string
}

type TemplateInfo struct {
	Name                  string
	DefaultOutputFileName string
	Description           string
	Extension             string
	TemplateFileName      string
}

/*type Language struct {
	Name            string
	DefaultFileName string
	Description     string
	Extension       string
}*/
/*
func loadLanguageInfoFile(languageName string) Language {
	data, err := templatesDir.ReadFile(languageName)
	config := Language{}

	if err != nil {
		log.Fatal(errors.New("Failed to load config file: " + languageName))
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(errors.New("Failed to read config file: " + languageName))
	}

	return config
}

*/
func loadLanguageConfigFile(languageName string) LanguageConfig {
	data, err := templatesDir.ReadFile(languageName)
	config := LanguageConfig{}

	if err != nil {
		log.Fatal(errors.New("Failed to load config file: " + languageName))
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(errors.New("Failed to read config file: " + languageName))
	}

	return config
}
