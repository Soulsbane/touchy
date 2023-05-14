package scripts

import (
	"embed"
	"fmt"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/api"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/path"
	"github.com/Soulsbane/touchy/internal/templates"
	"os"
	"path/filepath"
)

const defaultScriptFileName = "main.lua"

//go:embed scripts
var embedsDir embed.FS

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
	scripts      map[string]infofile.InfoFile
}

func New(languageTemplates *templates.Templates) *TouchyScripts {
	var touchyScripts TouchyScripts

	touchyScripts.scripts = make(map[string]infofile.InfoFile)
	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewScriptErrors())
	touchyScripts.scriptSystem.SetGlobal("Templates", languageTemplates)
	touchyScripts.registerFunctions()
	touchyScripts.findEmbeddedScripts()
	touchyScripts.findConfigDirScripts()

	return &touchyScripts
}

func (ts *TouchyScripts) findConfigDirScripts() {
	configScriptsDir := path.GetScriptsDir()
	configDirs, err := os.ReadDir(configScriptsDir)

	if err != nil {
		fmt.Println("Failed to read config directory: ", err)
	}

	for _, dir := range configDirs {
		if dir.IsDir() {
			defaultConfig := infofile.InfoFile{
				Name:        dir.Name(),
				Description: "<Unknown>",
				Embedded:    false,
			}

			infoFileName := filepath.Join("scripts", dir.Name(), infofile.DefaultFileName)
			config, err := infofile.Load(infoFileName, embedsDir)

			if err != nil {
				ts.scripts[dir.Name()] = defaultConfig
			} else {
				ts.scripts[dir.Name()] = config
			}
		}
	}
}

func (ts *TouchyScripts) findEmbeddedScripts() {
	dirs, err := embedsDir.ReadDir("scripts")

	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			defaultConfig := infofile.InfoFile{
				Name:        dir.Name(),
				Description: "<Unknown>",
				Embedded:    true,
			}

			infoFileName := filepath.Join("scripts", dir.Name(), infofile.DefaultFileName)
			config, err := infofile.Load(infoFileName, embedsDir)

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
		if ts.scripts[scriptName].Embedded {
			scriptPath := filepath.Join("scripts", scriptName, defaultScriptFileName)
			data, err := embedsDir.ReadFile(scriptPath)

			if err != nil {
				fmt.Println("Failed to load script: " + scriptName)
			} else {
				ts.scriptSystem.DoString(string(data))
			}
		} else {
			scriptPath := filepath.Join(path.GetScriptsDir(), scriptName, defaultScriptFileName)
			data, err := os.ReadFile(scriptPath)

			if err != nil {
				fmt.Println(err)
			} else {
				ts.scriptSystem.DoString(string(data))
			}
		}
	}
}
