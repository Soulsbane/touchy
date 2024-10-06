package main

import (
	"errors"
	"fmt"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/scripts"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/Soulsbane/touchy/internal/templates"
)

var languageTemplates *templates.Templates
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
	languageTemplates, userTemplatesErr, embeddedTemplatesErr = templates.New()

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

func handleCreateCommand(languageName string, templateName string, fileName string) {
	err := languageTemplates.CreateFileFromTemplate(languageName, templateName, fileName)

	if err != nil {
		handleError(err, templateName, languageName)
	}
}

func handleListCommand(listType string, languageName string) {
	switch listType {
	case "all":
		scriptsList := touchyScripts.GetListOfScripts()
		ListScripts(scriptsList)
		fmt.Println("")
		ListTemplates(languageName, languageTemplates.GetListOfAllLanguages())
	case "languageTemplates":
		ListLanguages(languageTemplates.GetListOfAllLanguages())
	case "scripts":
		scriptsList := touchyScripts.GetListOfScripts()
		ListScripts(scriptsList)
	case "templates":
		ListTemplates(languageName, languageTemplates.GetListOfAllLanguages())
	default:
		fmt.Println("That list type could not be found! Use 'list all' to see all available types.")
	}
}

func handleShowCommand(languageName string, templateName string) {
	err := languageTemplates.ShowTemplate(languageName, templateName)

	if err != nil {
		handleError(err, templateName, languageName)
	}
}

func handleRunCommand(scriptName string) {
	err := touchyScripts.Run(scriptName)

	if err != nil {
		handleError(err, scriptName, "")
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
			err := languageTemplates.CreateFileFromTemplate(createCmd.Language, createCmd.TemplateName, createCmd.FileName)

			if err != nil {
				handleError(err, createCmd.TemplateName, createCmd.Language)
			}
		}
	}
}
