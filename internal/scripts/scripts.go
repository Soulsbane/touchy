package scripts

import (
	"embed"
	"fmt"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/api"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/path"
	"github.com/Soulsbane/touchy/internal/templates"
	"github.com/Soulsbane/touchy/internal/ui"
	"github.com/jedib0t/go-pretty/v6/table"
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

func New() *TouchyScripts {
	var touchyScripts TouchyScripts

	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewStdOutScriptErrors())
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

	ts.findScripts(configDirs, false)
}

func (ts *TouchyScripts) findEmbeddedScripts() {
	dirs, err := embedsDir.ReadDir("scripts")

	if err != nil {
		fmt.Println("Failed to read embeds directory: ", err)
	}

	ts.findScripts(dirs, true)
}

func (ts *TouchyScripts) findScripts(dirs []fs.DirEntry, embedded bool) {
	for _, dir := range dirs {
		if dir.IsDir() {
			var data []byte
			var infoFilePath string
			var err error

			if embedded {
				infoFilePath = filepath.Join("scripts", dir.Name(), infofile.DefaultFileName)
				data, err = embedsDir.ReadFile(infoFilePath)

				if err != nil {
					fmt.Println("Failed to load config file: " + infoFilePath)
				}
			} else {
				infoFilePath = filepath.Join(path.GetScriptsDir(), dir.Name(), infofile.DefaultFileName)
				data, err = os.ReadFile(infoFilePath)

				if err != nil {
					fmt.Println("Failed to load config file: " + infoFilePath)
				}
			}

			config := infofile.Load(dir.Name(), infoFilePath, embedded, data)
			ts.scripts = append(ts.scripts, config)
		}
	}
}

func (ts *TouchyScripts) List(listArg string) {
	outputTable := ui.CreateNewTableWriter("Scripts", "Script Name", "Description")

	for _, script := range ts.scripts {
		outputTable.AppendRow(table.Row{script.Name, script.Description})
	}

	outputTable.Render()
}

func (ts *TouchyScripts) RegisterAPI() {
	ts.scriptSystem.SetGlobal("GetOutputDir", api.GetOutputDir)
	ts.scriptSystem.SetGlobal("GetAppConfigDir", path.GetAppConfigDir)
	ts.scriptSystem.SetGlobal("GetScriptsDir", path.GetScriptsDir)
	ts.scriptSystem.SetGlobal("GetTemplatesDir", path.GetTemplatesDir)

	templatesObject := templates.New()
	pathAPI := &api.Path{}

	ts.scriptSystem.SetGlobal("Templates", templatesObject)
	ts.scriptSystem.SetGlobal("Path", pathAPI)
}

func (ts *TouchyScripts) Run(scriptName string) error {
	idx := slices.IndexFunc(ts.scripts, func(c infofile.InfoFile) bool { return c.Name == scriptName })

	if idx >= 0 {
		scriptInfo := ts.scripts[idx]
		if scriptInfo.Name == scriptName {
			if scriptInfo.Embedded {
				scriptPath := filepath.Join("scripts", scriptName, defaultScriptFileName)
				data, err := embedsDir.ReadFile(scriptPath)

				if err != nil {
					return fmt.Errorf("failed to load script: %s", scriptName)
				} else {
					ts.scriptSystem.DoString(string(data))
					return nil
				}
			} else {
				scriptPath := filepath.Join(path.GetScriptsDir(), scriptName, defaultScriptFileName)
				data, err := os.ReadFile(scriptPath)

				if err != nil {
					return fmt.Errorf("failed to load script: %s", scriptName)
				} else {
					ts.scriptSystem.DoString(string(data))
					return nil
				}
			}
		}
	} else {
		return fmt.Errorf("script not found: %s", scriptName)
	}

	return nil
}
