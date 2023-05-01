package scripts

import (
	"embed"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/infofile"
	"path/filepath"
)

//go:embed scripts
var scriptsDir embed.FS

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
	scripts      map[string]infofile.InfoFile
}

func New() *TouchyScripts {
	var touchyScripts TouchyScripts

	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewScriptErrors())
	touchyScripts.registerFunctions()
	touchyScripts.findScripts()

	return &touchyScripts
}

func (ts *TouchyScripts) findScripts() {
	dirs, err := scriptsDir.ReadDir("scripts")

	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			defaultConfig := infofile.InfoFile{
				Name:        dir.Name(),
				Description: "<Unknown>",
			}

			ts.scripts = make(map[string]infofile.InfoFile)
			infoFileName := filepath.Join("scripts", dir.Name(), infofile.DefaultFileName)
			config, err := infofile.Load(infoFileName, scriptsDir)

			if err != nil {
				ts.scripts[dir.Name()] = defaultConfig
			} else {
				ts.scripts[dir.Name()] = config
			}
		}
	}
}

func (ts *TouchyScripts) registerFunctions() {

}

func (ts *TouchyScripts) Run(languageName string, scriptName string) {

}
