package main

import (
	"errors"
	"fmt"
	"github.com/Soulsbane/touchy/internal/scripts"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/templates"
)

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

func handleCreateCommand(languages *templates.Templates, languageName string, templateName string, fileName string) {
	err := languages.CreateFileFromTemplate(languageName, templateName, fileName)

	if err != nil {
		handleError(err, templateName, languageName)
	}
}

func handleListCommand(languages *templates.Templates, listType string, languageName string) {
	scriptToRun := scripts.New()

	scriptToRun.RegisterAPI()

	switch listType {
	case "all":
		scriptsList := scriptToRun.GetListOfScripts()
		ListScripts(scriptsList)
		fmt.Println("")
		ListTemplates(languageName, languages.GetListOfAllLanguages())
	case "languages":
		ListLanguages(languages.GetListOfAllLanguages())
	case "scripts":
		scriptsList := scriptToRun.GetListOfScripts()
		ListScripts(scriptsList)
	case "templates":
		ListTemplates(languageName, languages.GetListOfAllLanguages())
	default:
		fmt.Println("That list type could not be found! Use 'list all' to see all available types.")
	}
}

func main() {
	var cmds commands

	cmdLineArgs := os.Args[1:]
	cmd := cmdLineArgs[0]

	err := pathutils.SetupConfigDir()

	if err != nil {
		fmt.Println("Failed to setup config directory: ", err)
	}

	if len(cmdLineArgs) == 0 {
		fmt.Println("No arguments provided. Use -h or --help for more information.")
	} else {
		languages, userTemplatesErr, embeddedTemplatesErr := templates.New()

		if userTemplatesErr != nil || embeddedTemplatesErr != nil {
			handleError(userTemplatesErr, "", "")
			handleError(embeddedTemplatesErr, "", "")
		}

		if isReservedCommand(cmds, cmd) || cmd == "-h" || cmd == "--help" {
			arg.MustParse(&cmds)

			switch {
			case cmds.Create != nil:
				handleCreateCommand(languages, cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)
			case cmds.List != nil:
				handleListCommand(languages, cmds.List.Type, cmds.List.Language)
			case cmds.Show != nil:
				err := languages.ShowTemplate(cmds.Show.Language, cmds.Show.TemplateName)

				if err != nil {
					handleError(err, cmds.Show.TemplateName, cmds.Show.Language)
				}
			case cmds.Run != nil:
				scriptToRun := scripts.New()

				scriptToRun.RegisterAPI()
				err := scriptToRun.Run(cmds.Run.ScriptName)

				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			var createCmd CreateCommand

			arg.MustParse(&createCmd)
			err := languages.CreateFileFromTemplate(createCmd.Language, createCmd.TemplateName, createCmd.FileName)

			if err != nil {
				handleError(err, createCmd.TemplateName, createCmd.Language)
			}
		}
	}
}
