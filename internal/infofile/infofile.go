package infofile

import (
	"github.com/pelletier/go-toml/v2"
)

const DefaultFileName = "info.toml"
const DefaultDescription = "No description available"
const defaultOutputFileName = "output.txt"
const defaultName = "No name provided"

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
		Description: DefaultDescription,
		Embedded:    embedded,
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		// INFO: Really we don't need to worry about the error here since default values are fine
		return GetDefaultInfoFile()
	}

	return config
}

func GetDefaultInfoFile() InfoFile {
	return InfoFile{
		Name:                  defaultName,
		DefaultOutputFileName: defaultOutputFileName,
		Description:           DefaultDescription,
		Embedded:              false,
	}
}
