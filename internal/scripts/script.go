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
	Run(scriptName string, system goscriptsystem.ScriptSystem) error
}

type TouchyScript struct {
	scriptSystem *goscriptsystem.ScriptSystem
	info         infofile.InfoFile
}
