package scripts

import (
	"github.com/Soulsbane/goscriptsystem/goscriptsystem"
	"github.com/Soulsbane/touchy/internal/infofile"
)

type TouchyScript struct {
	scriptSystem *goscriptsystem.ScriptSystem
	info         infofile.InfoFile
}
