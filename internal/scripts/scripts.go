package scripts

import (
	"embed"
	"fmt"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/api"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/path"
	"github.com/Soulsbane/touchy/internal/templates"
	"golang.org/x/exp/slices"
	"io/fs"
	"os"
	"path/filepath"
)

const defaultScriptFileName = "main.lua"

//go:embed scripts
var embedsDir embed.FS

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
	scripts      []infofile.InfoFile
}

func New(languageTemplates *templates.Templates) *TouchyScripts {
	var touchyScripts TouchyScripts

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

	ts.findScripts(configDirs, false, embedsDir)
}

func (ts *TouchyScripts) findEmbeddedScripts() {
	dirs, err := embedsDir.ReadDir("scripts")

	if err != nil {
		fmt.Println("Failed to read embeds directory: ", err)
	}

	ts.findScripts(dirs, true, embedsDir)
}

func (ts *TouchyScripts) findScripts(dirs []fs.DirEntry, embedded bool, fs embed.FS) {
	for _, dir := range dirs {
		if dir.IsDir() {
			defaultConfig := infofile.InfoFile{
				Name:        dir.Name(),
				Description: "<Unknown>",
				Embedded:    embedded,
			}

			infoFileName := filepath.Join("scripts", dir.Name(), infofile.DefaultFileName)
			config, err := infofile.Load(infoFileName, fs)

			if err != nil {
				ts.scripts = append(ts.scripts, defaultConfig)
			} else {
				config.Embedded = embedded
				ts.scripts = append(ts.scripts, config)
			}
		}
	}
}

func (ts *TouchyScripts) registerFunctions() {
	ts.scriptSystem.SetGlobal("GetOutputDir", api.GetOutputDir)
	ts.scriptSystem.SetGlobal("GetAppConfigDir", path.GetAppConfigDir)
}

func (ts *TouchyScripts) Run(scriptName string) {
	idx := slices.IndexFunc(ts.scripts, func(c infofile.InfoFile) bool { return c.Name == scriptName })

	if idx >= 0 {
		scriptInfo := ts.scripts[idx]
		if scriptInfo.Name == scriptName {
			if scriptInfo.Embedded {
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
	} else {
		fmt.Println("Script not found: " + scriptName)
	}
}
