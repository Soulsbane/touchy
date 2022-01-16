package generator

import (
	"errors"
	"log"

	"github.com/pelletier/go-toml/v2"
)

type Language struct {
	Name            string
	DefaultFileName string
	Description     string
	Extension       string
}

func LoadLanguageConfigFile(languageName string) Language {
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
