package scripts

import (
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/infofile"
)

type Scripts interface {
	GetScriptInfoFor(scriptName string) infofile.InfoFile
	GetListOfScripts() []TouchyScript
}

type TouchyScript struct {
	scriptSystem *goscriptsystem.ScriptSystem
	info         infofile.InfoFile
}
