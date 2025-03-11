package scripts

import (
	"fmt"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/infofile"
	"golang.org/x/exp/slices"
	"path"
)

type EmbeddedScripts struct {
	scripts []TouchyScript
}

func NewEmbeddedScripts() (*EmbeddedScripts, error) {
	var embeddedScripts EmbeddedScripts
	err := embeddedScripts.findScripts()

	if err != nil {
		return &embeddedScripts, err
	}

	return &embeddedScripts, nil
}

func (es *EmbeddedScripts) findScripts() error {
	dirs, err := embedsDir.ReadDir("scripts")

	if err != nil {
		return fmt.Errorf("failed to read script embeds directory: %w", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			var touchyScript TouchyScript

			touchyScript.scriptSystem = goscriptsystem.New(goscriptsystem.NewStdOutScriptErrors())
			infoFilePath := path.Join("scripts", dir.Name(), infofile.DefaultFileName)
			data, readFileErr := embedsDir.ReadFile(infoFilePath)

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

func (es *EmbeddedScripts) GetListOfScripts() []TouchyScript {
	return es.scripts
}

func (es *EmbeddedScripts) GetListOfScriptInfo() []infofile.InfoFile {
	var infoList []infofile.InfoFile

	for _, script := range es.scripts {
		infoList = append(infoList, script.info)
	}

	return infoList
}

func (es *EmbeddedScripts) GetScriptInfoFor(scriptName string) infofile.InfoFile {
	idx := slices.IndexFunc(es.scripts, func(c TouchyScript) bool { return c.info.GetName() == scriptName })

	if idx >= 0 {
		return es.scripts[idx].info
	}

	return infofile.InfoFile{}
}
