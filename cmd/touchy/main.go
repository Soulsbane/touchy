package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/scripts"
	"github.com/Soulsbane/touchy/internal/templates"
)

func handleError(err error, templateName string, languageName string) {
	if err != nil {
		if errors.Is(err, templates.ErrLanguageNotFound) {
			fmt.Println("Language not found:", languageName)
		} else if errors.Is(err, templates.ErrTemplateNotFound) {
			fmt.Println("Template not found:", templateName)
		} else if errors.Is(err, templates.ErrFileNameEmpty) {
			fmt.Println("Error: output filename not specified")
		} else {
			fmt.Println("Error:", err)
		}
	}
}

func main() {
	var cmds commands
	cmdLineArgs := os.Args[1:]

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

		scriptToRun := scripts.New()
		cmd := cmdLineArgs[0]

		scriptToRun.RegisterAPI()

		if isReservedCommand(cmds, cmd) || cmd == "-h" || cmd == "--help" {
			arg.MustParse(&cmds)

			switch {
			case cmds.Create != nil:
				err := languages.CreateFileFromTemplate(cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)

				if err != nil {
					handleError(err, cmds.Create.TemplateName, cmds.Create.Language)
				}
			case cmds.List != nil:
				switch cmds.List.Type {
				case "all":
					scriptsList := scriptToRun.GetListOfScripts()
					ListScripts(scriptsList)
					fmt.Println("")
					ListTemplates(cmds.List.Language, languages.GetListOfAllLanguages())
				case "languages":
					ListLanguages(languages.GetListOfAllLanguages())
				case "scripts":
					scriptsList := scriptToRun.GetListOfScripts()
					ListScripts(scriptsList)
				case "templates":
					ListTemplates(cmds.List.Language, languages.GetListOfAllLanguages())
				}
			case cmds.Show != nil:
				err := languages.ShowTemplate(cmds.Show.Language, cmds.Show.TemplateName)

				if err != nil {
					handleError(err, cmds.Show.TemplateName, cmds.Show.Language)
				}
			case cmds.Run != nil:
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
