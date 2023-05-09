package scripts

import (
	"embed"
	"fmt"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/api"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/path"
	"github.com/Soulsbane/touchy/internal/templates"
	"path/filepath"
)

const defaultScriptFileName = "main.lua"

//go:embed scripts
var scriptsDir embed.FS

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
	scripts      map[string]infofile.InfoFile
}

func New(languageTemplates *templates.Templates) *TouchyScripts {
	var touchyScripts TouchyScripts

	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewScriptErrors())
	touchyScripts.scriptSystem.SetGlobal("Templates", languageTemplates)
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
	ts.scriptSystem.SetGlobal("GetOutputDir", api.GetOutputDir)
	ts.scriptSystem.SetGlobal("GetAppConfigDir", path.GetAppConfigDir)
}

func (ts *TouchyScripts) Run(scriptName string) {
	if _, ok := ts.scripts[scriptName]; !ok {
		fmt.Println("Script not found: " + scriptName)
	} else {
		scriptPath := filepath.Join("scripts", scriptName, defaultScriptFileName)
		data, err := scriptsDir.ReadFile(scriptPath)

		if err != nil {
			fmt.Println("Failed to load script: " + scriptName)
		} else {
			ts.scriptSystem.DoString(string(data))
		}
	}
}
