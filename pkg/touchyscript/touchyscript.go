package touchyscript

import "github.com/Soulsbane/goscriptsystem/goscriptsystem"

type TouchyScript struct {
	scriptSystem *goscriptsystem.ScriptSystem
}

func New(scriptName string) *TouchyScript {
	var touchyScript TouchyScript

	touchyScript.scriptSystem = goscriptsystem.New(goscriptsystem.NewScriptErrors())
	return &touchyScript
}
