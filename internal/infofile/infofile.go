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

func (i *InfoFile) GetName() string {
	return i.Name
}

func (i *InfoFile) GetDefaultOutputFileName() string {
	return i.DefaultOutputFileName
}

func (i *InfoFile) GetDescription() string {
	return i.Description
}

func (i *InfoFile) IsEmbedded() bool {
	return i.Embedded
}

func (i *InfoFile) SetEmbedded(embedded bool) {
	i.Embedded = embedded
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
