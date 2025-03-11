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
var scriptsManager *scripts.TouchyScriptsManager

func setupScriptsAndTemplates() {
	pathUtilsErr := pathutils.SetupConfigDir()

	if pathUtilsErr != nil {
		fmt.Println("Failed to setup config directory: ", pathUtilsErr)
	}

	//touchyScripts.RegisterAPI()

	manager = templates.New()
	scriptsManager = scripts.New()
	manager.GatherTemplates()
	scriptsManager.GatherScripts()
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

func main() {
	var cmds commands
	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) == 0 {
		fmt.Println("No arguments provided. Use -h or --help for more information.")
	} else {
		cmd := cmdLineArgs[0]
		setupScriptsAndTemplates()

		if isReservedCommand(cmds, cmd) || cmd == "-h" || cmd == "--help" {
			arg.MustParse(&cmds)

			switch {
			case cmds.Create != nil:
				handleCreateCommand(cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)
			case cmds.List != nil:
				handleListCommand(cmds.List.Type, cmds.List.Language)
			case cmds.Show != nil:
				handleShowCommand(cmds.Show.Language, cmds.Show.TemplateName)
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
