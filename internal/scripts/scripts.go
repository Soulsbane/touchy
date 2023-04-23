package scripts

import (
	"embed"
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
)

//go:embed scripts
var scriptsDir embed.FS

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
}

func New() *TouchyScripts {
	var touchyScripts TouchyScripts

	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewScriptErrors())
	return &touchyScripts
}

func (ts *TouchyScripts) Run(languageName string, scriptName string) {

}
