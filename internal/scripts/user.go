package scripts

import (
	"fmt"
	"os"
	"path"

	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"golang.org/x/exp/slices"
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

func (us *UserScripts) findScripts() error {
	scriptsPath := pathutils.GetScriptsDir()
	dirs, err := os.ReadDir(scriptsPath)

	if err != nil {
		return fmt.Errorf("failed to read script embeds directory: %w", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			var touchyScript TouchyScript
			var config infofile.InfoFile

			infoFilePath := path.Join(scriptsPath, dir.Name(), infofile.DefaultFileName)
			data, readFileErr := os.ReadFile(infoFilePath)

			if readFileErr != nil {
				config = infofile.GetDefaultInfoFile()
			}

			config = infofile.Load(dir.Name(), infoFilePath, true, data)
			config.SetEmbedded(true)
			touchyScript.info = config
			us.scripts = append(us.scripts, touchyScript)
		}
	}

	return nil
}

func (us *UserScripts) GetListOfScripts() []TouchyScript {
	return us.scripts
}

func (us *UserScripts) GetListOfScriptInfo() []infofile.InfoFile {
	var infoList []infofile.InfoFile

	for _, script := range us.scripts {
		infoList = append(infoList, script.info)
	}

	return infoList
}

func (us *UserScripts) GetScriptInfoFor(scriptName string) infofile.InfoFile {
	idx := slices.IndexFunc(us.scripts, func(c TouchyScript) bool { return c.info.GetName() == scriptName })

	if idx >= 0 {
		return us.scripts[idx].info
	}

	return infofile.InfoFile{}
}

func (us *UserScripts) Run(scriptName string, system goscriptsystem.ScriptSystem) error {
	idx := slices.IndexFunc(us.scripts, func(c TouchyScript) bool { return c.info.GetName() == scriptName })
	scriptsPath := pathutils.GetScriptsDir()

	if idx >= 0 {
		script := us.scripts[idx]
		script.scriptSystem = &system
		scriptPath := path.Join(scriptsPath, scriptName, defaultScriptFileName)
		data, err := os.ReadFile(scriptPath)

		if err != nil {
			return fmt.Errorf("failed to read script file: %w", err)
		} else {
			script.scriptSystem.DoString(string(data))
			return nil
		}
	}

	return fmt.Errorf("script %s not found", scriptName)
}
