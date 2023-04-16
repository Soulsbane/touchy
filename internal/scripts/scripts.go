package scripts

import "github.com/Soulsbane/goscriptsystem/goscriptsystem"

type TouchyScripts struct {
	scriptSystem *goscriptsystem.ScriptSystem
}

func New() *TouchyScripts {
	var touchyScripts TouchyScripts

	touchyScripts.scriptSystem = goscriptsystem.New(goscriptsystem.NewScriptErrors())
	return &touchyScripts
}
