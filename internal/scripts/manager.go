package scripts

type TouchyScriptsManager struct {
	scripts []Scripts
}

func New() *TouchyScriptsManager {
	var manager TouchyScriptsManager
	return &manager
}

func (manager *TouchyScriptsManager) GatherScripts() {
	embedded, embeddedErr := NewEmbeddedScripts()

	if embeddedErr != nil {
		panic(embeddedErr)
	}

	manager.scripts = append(manager.scripts, embedded)
}

func (manager *TouchyScriptsManager) GetListOfScripts() []TouchyScript {
	var scriptList []TouchyScript

	for _, script := range manager.scripts {
		scriptList = append(scriptList, script.GetListOfScripts()...)
	}

	return scriptList
}
