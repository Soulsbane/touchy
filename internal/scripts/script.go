package scripts

import "github.com/Soulsbane/goscriptsystem/goscriptsystem"

type TouchyScript struct {
	scriptSystem *goscriptsystem.ScriptSystem
	Name         string
	Author       string
	Description  string
	Embedded     bool
}

func NewTouchyScript() *TouchyScript {
	var touchyScript TouchyScript

	touchyScript.scriptSystem = goscriptsystem.New(goscriptsystem.NewStdOutScriptErrors())
	return &touchyScript
}
