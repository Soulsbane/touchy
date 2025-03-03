package scripts

import "github.com/Soulsbane/goscriptsystem/goscriptsystem"

type TouchyScript struct {
	scriptSystem *goscriptsystem.ScriptSystem
	Name         string
	Author       string
	Description  string
	Embedded     bool
}
