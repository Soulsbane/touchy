package scripts

import (
	"fmt"

	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/api"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/templates"
	libs "github.com/vadv/gopher-lua-libs"
)

type TouchyScriptsManager struct {
	scripts []Scripts
}

func New() *TouchyScriptsManager {
	var manager TouchyScriptsManager
	return &manager
}

func (manager *TouchyScriptsManager) createScriptSystem() *goscriptsystem.ScriptSystem {
	scriptSystem := goscriptsystem.New(goscriptsystem.NewStdOutScriptErrors())

	scriptSystem.SetGlobal("GetOutputDir", pathutils.GetOutputDir)
	scriptSystem.SetGlobal("GetAppConfigDir", pathutils.GetAppConfigDir)
	scriptSystem.SetGlobal("GetScriptsDir", pathutils.GetScriptsDir)
	scriptSystem.SetGlobal("GetTemplatesDir", pathutils.GetTemplatesDir)
	scriptSystem.SetGlobal("CleanPath", pathutils.CleanPath)

	templatesObject := templates.New()
	templatesObject.GatherTemplates()
	promptsObject := api.NewPrompts()
	ioObject := api.NewIO()
	commandObject := api.NewCommand()
	downloaderObject := api.NewDownloader()

	scriptSystem.SetGlobal("Templates", templatesObject)
	scriptSystem.SetGlobal("Prompts", promptsObject)
	scriptSystem.SetGlobal("IO", ioObject)
	scriptSystem.SetGlobal("Command", commandObject)
	scriptSystem.SetGlobal("Downloader", downloaderObject)
	libs.Preload(scriptSystem.GetState())

	return scriptSystem
}

func (manager *TouchyScriptsManager) GatherScripts() {
	embedded, embeddedErr := NewEmbeddedScripts()
	user, userErr := NewUserScripts()

	if embeddedErr != nil {
		panic(embeddedErr)
	}

	if userErr != nil {
		panic(userErr)
	}

	manager.scripts = append(manager.scripts, embedded)
	manager.scripts = append(manager.scripts, user)
}

func (manager *TouchyScriptsManager) GetListOfScripts() []TouchyScript {
	var scriptList []TouchyScript

	for _, script := range manager.scripts {
		scriptList = append(scriptList, script.GetListOfScripts()...)
	}

	return scriptList
}

func (manager *TouchyScriptsManager) GetListOfScriptInfo() []infofile.InfoFile {
	var scriptList []infofile.InfoFile

	for _, script := range manager.scripts {
		scriptList = append(scriptList, script.GetListOfScriptInfo()...)
	}

	return scriptList
}

func (manager *TouchyScriptsManager) GetNumberOfScripts(name string) int {
	count := 0

	for _, script := range manager.scripts {
		if script.HasScript(name) {
			count++
		}
	}

	return count
}

func (manager *TouchyScriptsManager) Run(scriptName string, args []string) error {
	scriptSystem := manager.createScriptSystem()

	for _, script := range manager.scripts {
		if script.HasScript(scriptName) {

			err := script.Run(scriptName, *scriptSystem, args)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("script %s not found", scriptName)
}
