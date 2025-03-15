package scripts

import (
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
	scriptSystem.SetGlobal("DownloadFile", api.DownloadFile)
	scriptSystem.SetGlobal("DownloadFileWithProgress", api.DownloadFileWithProgress)

	templatesObject := templates.New()
	templatesObject.GatherTemplates()
	promptsObject := api.NewPrompts()
	ioObject := api.NewIO()

	scriptSystem.SetGlobal("Templates", templatesObject)
	scriptSystem.SetGlobal("Prompts", promptsObject)
	scriptSystem.SetGlobal("IO", ioObject)
	libs.Preload(scriptSystem.GetState())

	return scriptSystem
}

func (manager *TouchyScriptsManager) GatherScripts() {
	embedded, embeddedErr := NewEmbeddedScripts()

	if embeddedErr != nil {
		panic(embeddedErr)
	}

	manager.scripts = append(manager.scripts, embedded)
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
func (manager *TouchyScriptsManager) Run(scriptName string) error {
	scriptSystem := manager.createScriptSystem()

	for _, script := range manager.scripts {
		err := script.Run(scriptName, *scriptSystem)

		if err != nil {
			return err
		}
	}

	return nil
}
