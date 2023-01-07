package templates

import (
	"errors"

	"github.com/pelletier/go-toml/v2"
)

type Language struct {
	// dirName         string                  // Name of the directory under the templates directory.
	infoConfig      CommonConfig            // Each language has a config file in its root directory call config.toml
	templateConfigs map[string]CommonConfig // A list of all the templates in the language directory.
}

type CommonConfig struct {
	Name                  string
	DefaultOutputFileName string
	Description           string
}

func loadInfoFile(languageName string) (CommonConfig, error) {
	data, err := templatesDir.ReadFile(languageName)
	config := CommonConfig{}

	if err != nil {
		return config, errors.New("Failed to load config file: " + languageName)
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		return config, errors.New("Failed to read config data: " + languageName)
	}

	return config, nil
}
