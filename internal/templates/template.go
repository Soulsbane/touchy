package templates

import (
	"embed"
	"errors"

	"github.com/pelletier/go-toml/v2"
)

type Language struct {
	// dirName         string                  // Name of the directory under the templates directory.
	infoConfig      CommonConfig            // Each language has a config file in its root directory call config.toml
	templateConfigs map[string]CommonConfig // A list of all the templates in the language directory. The key is the template dir name.
}

type CommonConfig struct {
	Name                  string
	DefaultOutputFileName string
	Description           string
}

// TODO: Refactor this info file loading code into its own package
func LoadInfoFile(languageName string, fs embed.FS) (CommonConfig, error) {
	data, err := fs.ReadFile(languageName)
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
