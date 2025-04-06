package scripts

import (
	"fmt"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"golang.org/x/exp/slices"
	"os"
	"path"
)

type UserScripts struct {
	scripts []TouchyScript
}

func NewUserScripts() (*UserScripts, error) {
	var userScripts UserScripts
	err := userScripts.findScripts()

	if err != nil {
		return &userScripts, err
	}

	return &userScripts, nil
}

func (es *UserScripts) findScripts() error {
	scriptsPath := pathutils.GetScriptsDir()
	dirs, err := os.ReadDir(scriptsPath)

	if err != nil {
		return fmt.Errorf("failed to read script embeds directory: %w", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			var touchyScript TouchyScript

			infoFilePath := path.Join(scriptsPath, dir.Name(), infofile.DefaultFileName)
			data, readFileErr := os.ReadFile(infoFilePath)

			if readFileErr != nil {
				// TODO: Maybe set a default config if config file is not found?
				fmt.Println("Failed to load config file: " + infoFilePath)
			}

			config := infofile.Load(dir.Name(), infoFilePath, true, data)
			config.SetEmbedded(true)
			touchyScript.info = config
			es.scripts = append(es.scripts, touchyScript)
		}
	}

	return nil
}

func (es *UserScripts) GetListOfScripts() []TouchyScript {
	return es.scripts
}

func (es *UserScripts) GetListOfScriptInfo() []infofile.InfoFile {
	var infoList []infofile.InfoFile

	for _, script := range es.scripts {
		infoList = append(infoList, script.info)
	}

	return infoList
}

func (es *UserScripts) GetScriptInfoFor(scriptName string) infofile.InfoFile {
	idx := slices.IndexFunc(es.scripts, func(c TouchyScript) bool { return c.info.GetName() == scriptName })

	if idx >= 0 {
		return es.scripts[idx].info
	}

	return infofile.InfoFile{}
}

func (es *UserScripts) Run(scriptName string, system goscriptsystem.ScriptSystem) error {
	idx := slices.IndexFunc(es.scripts, func(c TouchyScript) bool { return c.info.GetName() == scriptName })
	scriptsPath := pathutils.GetScriptsDir()

	if idx >= 0 {
		script := es.scripts[idx]
		script.scriptSystem = &system
		scriptPath := path.Join(scriptsPath, scriptName, defaultScriptFileName)
		data, err := os.ReadFile(scriptPath)

		if err != nil {
			return fmt.Errorf("failed to read script file: %w", err)
		} else {
			script.scriptSystem.DoString(string(data))
			return nil
		}
	} else {
		return fmt.Errorf("script %s not found", scriptName)
	}
}
