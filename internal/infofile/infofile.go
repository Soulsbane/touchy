package infofile

import (
	"embed"
	"errors"
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
