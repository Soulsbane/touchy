package infofile

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
)

const DefaultFileName = "info.toml"

type InfoFile struct {
	Name                  string
	DefaultOutputFileName string
	Description           string
	Embedded              bool
}

func Load(name string, infoFilePath string, embedded bool, data []byte) InfoFile {
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
