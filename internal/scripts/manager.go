package scripts

type TouchyScriptsManager struct {
	scripts []TouchyScript
}

func New() *TouchyScriptsManager {
	var manager TouchyScriptsManager
	return &manager
}
