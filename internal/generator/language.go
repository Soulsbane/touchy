package generator

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Language struct {
	Name            string
	DefaultFileName string
	Description     string
	Extension       string
}

func LoadLanguageConfigFile(languageName string) Language {
	exePath, _ := os.Executable()
	configFileName := filepath.Join(filepath.Dir(exePath), "../../internal/generator/templates/", languageName, "/config.toml")

	data, err := ioutil.ReadFile(configFileName)
	config := Language{}

	if err != nil {
		log.Fatal(errors.New("Failed to load config file: " + configFileName))
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(errors.New("Failed to read config file: " + configFileName))
	}

	return config
}
