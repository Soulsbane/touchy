package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/scripts"
	"github.com/Soulsbane/touchy/internal/templates"
)

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
		languages := templates.New()
		scriptToRun := scripts.New()
		cmd := cmdLineArgs[0]

		scriptToRun.RegisterAPI()

		if isReservedCommand(cmds, cmd) || cmd == "-h" || cmd == "--help" {
			arg.MustParse(&cmds)

			switch {
			case cmds.Create != nil:
				err := languages.CreateFileFromTemplate(cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)

				if err != nil {
					fmt.Println(err)
				}
			case cmds.List != nil:
				switch cmds.List.Type {
				case "all":
					scriptsList := scriptToRun.GetListOfScripts()
					ListScripts(scriptsList)
					fmt.Println("")
					ListTemplates(cmds.List.Language)
				case "languages":
					ListLanguages()
				case "scripts":
					scriptsList := scriptToRun.GetListOfScripts()
					ListScripts(scriptsList)
				case "templates":
					ListTemplates(cmds.List.Language)
				}
			case cmds.Show != nil:
				err := languages.ShowTemplate(cmds.Show.Language, cmds.Show.TemplateName)

				if err != nil {
					fmt.Println(err)
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
				fmt.Println(err)
			}
		}
	}
}
