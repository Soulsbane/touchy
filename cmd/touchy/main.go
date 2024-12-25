package main

import (
	"errors"
	"fmt"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/scripts"
	"github.com/Soulsbane/touchy/internal/templates"
	"github.com/alexflint/go-arg"
	//"maps"
	"os"
)

var manager *templates.TemplateManager
var touchyScripts *scripts.TouchyScripts

func setupScriptsAndTemplates() {
	var userTemplatesErr error
	var embeddedTemplatesErr error

	pathUtilsErr := pathutils.SetupConfigDir()

	if pathUtilsErr != nil {
		fmt.Println("Failed to setup config directory: ", pathUtilsErr)
	}

	touchyScripts = scripts.New()
	touchyScripts.RegisterAPI()
	manager, userTemplatesErr, embeddedTemplatesErr = templates.New()

	if userTemplatesErr != nil || embeddedTemplatesErr != nil {
		handleError(userTemplatesErr, "", "")
		handleError(embeddedTemplatesErr, "", "")
	}
}

func handleError(err error, templateName string, languageName string) {
	if err != nil {
		switch {
		case errors.Is(err, templates.ErrLanguageNotFound):
			fmt.Println("Language not found:", languageName)
		case errors.Is(err, templates.ErrTemplateNotFound):
			fmt.Println("Template not found:", templateName)
		case errors.Is(err, templates.ErrFileNameEmpty):
			fmt.Println("Error: output filename not specified")
		case errors.Is(err, templates.ErrNoUserTemplatesDir):
			// TODO: Maybe notify user where the user templates is or if it should be created?
			fmt.Println("Warning: no user templates found!")
		default:
			fmt.Println("Error:", err)
		}
	}
}

//func handleCreateCommand(languageName string, templateName string, fileName string) {
//	err := manager.CreateFileFromTemplate(languageName, templateName, fileName)
//
//	if err != nil {
//		handleError(err, templateName, languageName)
//	}
//}
//
//func handleListCommand(listType string, languageName string) {
//	switch listType {
//	case "all":
//		scriptsList := touchyScripts.GetListOfScripts()
//		ListScripts(scriptsList)
//		fmt.Println("")
//		ListTemplates(languageName, manager.GetListOfAllLanguages())
//	case "languages":
//		ListLanguages(manager.GetListOfAllLanguages())
//	case "scripts":
//		scriptsList := touchyScripts.GetListOfScripts()
//		ListScripts(scriptsList)
//	case "templates":
//		ListTemplates(languageName, manager.GetListOfAllLanguages())
//	default:
//		fmt.Println("That list type could not be found! Use 'list all' to see all available types.")
//	}
//}

//func handleShowCommand(languageName string, templateName string) {
//	err := manager.ShowTemplate(languageName, templateName)
//
//	if err != nil {
//		handleError(err, templateName, languageName)
//	}
//}
//
//func handleRunCommand(scriptName string) {
//	err := touchyScripts.Run(scriptName)
//
//	if err != nil {
//		handleError(err, scriptName, "")
//	}
//}

func oldgatherTemplates() {
	embedded := templates.NewEmbeddedTemplates()
	embedded.GetListOfAllLanguages()

	user := templates.NewUserTemplates()
	user.GetListOfAllLanguages()
}

func gatherTemplates() (map[string]templates.Language, []templates.Templates) {
	embedded := templates.NewEmbeddedTemplates()
	languages := embedded.GetListOfAllLanguages()
	//user := templates.NewUserTemplates()
	langTemps := []templates.Templates{embedded}

	//maps.Copy(languages, user.GetListOfAllLanguages())

	return languages, langTemps
}

func main() {
	var cmds commands
	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) == 0 {
		fmt.Println("No arguments provided. Use -h or --help for more information.")
	} else {
		cmd := cmdLineArgs[0]
		setupScriptsAndTemplates()
		languages, langTemps := gatherTemplates()

		if isReservedCommand(cmds, cmd) || cmd == "-h" || cmd == "--help" {
			arg.MustParse(&cmds)

			switch {
			case cmds.Create != nil:
				handleCreateCommand(cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)
			case cmds.List != nil:
				handleListCommand(cmds.List.Type, cmds.List.Language, languages, langTemps)
			case cmds.Show != nil:
				handleShowCommand(cmds.Show.Language, cmds.Show.TemplateName, langTemps)
			case cmds.Run != nil:
				handleRunCommand(cmds.Run.ScriptName)
			}
		} else {
			var createCmd CreateCommand

			arg.MustParse(&createCmd)
			err := manager.CreateFileFromTemplate(createCmd.Language, createCmd.TemplateName, createCmd.FileName)

			if err != nil {
				handleError(err, createCmd.TemplateName, createCmd.Language)
			}
		}
	}
}
