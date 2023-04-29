package scripts

import (
	"embed"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/templates"
	"path/filepath"
)

//go:embed scripts
var scriptsDir embed.FS

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
	scripts      map[string]templates.CommonConfig
}

func New() *TouchyScripts {
	var touchyScripts TouchyScripts

	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewScriptErrors())
	touchyScripts.findScripts()
	return &touchyScripts
}

func (touchyScripts *TouchyScripts) findScripts() {
	dirs, err := scriptsDir.ReadDir("scripts")

	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			defaultConfig := templates.CommonConfig{
				Name:        dir.Name(),
				Description: "<Unknown>",
			}

			touchyScripts.scripts = make(map[string]templates.CommonConfig)
			infoFileName := filepath.Join("scripts", dir.Name(), "info.toml")
			config, err := templates.LoadInfoFile(infoFileName, scriptsDir)

			if err != nil {
				touchyScripts.scripts[dir.Name()] = defaultConfig
			} else {
				touchyScripts.scripts[dir.Name()] = config
			}
		}
	}
}

func (ts *TouchyScripts) Run(languageName string, scriptName string) {

}
