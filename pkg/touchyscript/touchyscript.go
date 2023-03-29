package touchyscript

import "github.com/Soulsbane/goscriptsystem/goscriptsystem"

type TouchyScript struct {
	goscriptsystem.ScriptSystem
}

func New(scriptName string) *TouchyScript {
	var touchyScript TouchyScript

	return &touchyScript
}
