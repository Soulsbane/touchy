package scripts

import (
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/infofile"
)

const defaultScriptFileName = "main.lua"

type Scripts interface {
	GetScriptInfoFor(scriptName string) infofile.InfoFile
	GetListOfScripts() []TouchyScript
	GetListOfScriptInfo() []infofile.InfoFile
	HasScript(scriptName string) bool
	Run(scriptName string, system goscriptsystem.ScriptSystem, argsl []string) error
}

type TouchyScript struct {
	scriptSystem *goscriptsystem.ScriptSystem
	info         infofile.InfoFile
}
