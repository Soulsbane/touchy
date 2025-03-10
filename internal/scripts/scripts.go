package scripts

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/api"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/templates"
	libs "github.com/vadv/gopher-lua-libs"
	"golang.org/x/exp/slices"
)

const defaultScriptFileName = "main.lua"

//go:embed scripts
var embedsDir embed.FS

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
	scripts      []infofile.InfoFile
}

//func New() *TouchyScripts {
//	var touchyScripts TouchyScripts
//
//	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewStdOutScriptErrors())
//	touchyScripts.findEmbeddedScripts()
//	touchyScripts.findConfigDirScripts()
//
//	return &touchyScripts
//}

func (ts *TouchyScripts) findConfigDirScripts() {
	configScriptsDir := pathutils.GetScriptsDir()
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
				infoFilePath = path.Join("scripts", dir.Name(), infofile.DefaultFileName)
				data, err = embedsDir.ReadFile(infoFilePath)

				if err != nil {
					fmt.Println("Failed to load config file: " + infoFilePath)
				}
			} else {
				infoFilePath = path.Join(pathutils.GetScriptsDir(), dir.Name(), infofile.DefaultFileName)
				data, err = os.ReadFile(infoFilePath)

				if err != nil {
					fmt.Println("Failed to load config file: " + infoFilePath)
				}
			}

			config := infofile.Load(dir.Name(), infoFilePath, embedded, data)
			config.SetEmbedded(embedded)
			ts.scripts = append(ts.scripts, config)
		}
	}
}

func (ts *TouchyScripts) GetScriptInfoFor(scriptName string) *infofile.InfoFile {
	idx := slices.IndexFunc(ts.scripts, func(c infofile.InfoFile) bool { return c.GetName() == scriptName })

	if idx >= 0 {
		return &ts.scripts[idx]
	}

	return nil
}

func (ts *TouchyScripts) GetListOfScripts() []infofile.InfoFile {
	scripts := ts.scripts
	return scripts
}

func (ts *TouchyScripts) RegisterAPI() {
	ts.scriptSystem.SetGlobal("GetOutputDir", pathutils.GetOutputDir)
	ts.scriptSystem.SetGlobal("GetAppConfigDir", pathutils.GetAppConfigDir)
	ts.scriptSystem.SetGlobal("GetScriptsDir", pathutils.GetScriptsDir)
	ts.scriptSystem.SetGlobal("GetTemplatesDir", pathutils.GetTemplatesDir)
	ts.scriptSystem.SetGlobal("CleanPath", pathutils.CleanPath)
	ts.scriptSystem.SetGlobal("DownloadFile", api.DownloadFile)
	ts.scriptSystem.SetGlobal("DownloadFileWithProgress", api.DownloadFileWithProgress)

	templatesObject := templates.New()
	promptsObject := api.NewPrompts()
	ioObject := api.NewIO()

	ts.scriptSystem.SetGlobal("Templates", templatesObject)
	ts.scriptSystem.SetGlobal("Prompts", promptsObject)
	ts.scriptSystem.SetGlobal("IO", ioObject)
	libs.Preload(ts.scriptSystem.GetState())
}

func (ts *TouchyScripts) Run(scriptName string) error {
	idx := slices.IndexFunc(ts.scripts, func(c infofile.InfoFile) bool { return c.GetName() == scriptName })

	if idx >= 0 {
		scriptInfo := ts.scripts[idx]
		if scriptInfo.GetName() == scriptName {
			if scriptInfo.IsEmbedded() {
				scriptPath := path.Join("scripts", scriptName, defaultScriptFileName)
				data, err := embedsDir.ReadFile(scriptPath)

				if err != nil {
					return fmt.Errorf("failed to load script: %s", scriptName)
				} else {
					ts.scriptSystem.DoString(string(data))
					return nil
				}
			} else {
				scriptPath := path.Join(pathutils.GetScriptsDir(), scriptName, defaultScriptFileName)
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
