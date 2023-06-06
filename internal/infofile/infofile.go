package infofile

import (
	"embed"
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
)

const DefaultFileName = "info.toml"

type InfoFile struct {
	Name                  string
	DefaultOutputFileName string
	Description           string
	Embedded              bool
}

func LoadSimple(name string, infoFilePath string, embedded bool, data []byte) InfoFile {
	var err error
	config := InfoFile{
		Name:        name,
		Description: "<Unknown>",
		Embedded:    embedded,
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		fmt.Println("Failed to read config data: "+infoFilePath, " using default config.")
	}

	return config
}

func Load(languageName string, embedded bool, fs embed.FS) (InfoFile, error) {
	var data []byte
	var err error

	if embedded {
		data, err = fs.ReadFile(languageName)
	} else {
		data, err = os.ReadFile(languageName)
	}

	config := InfoFile{}

	if err != nil {
		return config, errors.New("Failed to load config file: " + languageName)
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		return config, errors.New("Failed to read config data: " + languageName)
	}

	return config, nil
}
