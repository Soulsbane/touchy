package templates

import (
	"github.com/Soulsbane/touchy/internal/infofile"
)

type Language struct {
	// dirName         string                  // Name of the directory under the templates directory.
	infoConfig      infofile.InfoFile            // Each language has a config file in its root directory call config.toml
	templateConfigs map[string]infofile.InfoFile // A list of all the templates in the language directory. The key is the template dir name.
}
