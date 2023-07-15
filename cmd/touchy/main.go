package main

import (
	"fmt"
	"os"

	"github.com/Soulsbane/touchy/internal/path"
	"github.com/Soulsbane/touchy/internal/scripts"
	"github.com/Soulsbane/touchy/internal/templates"
	"github.com/alexflint/go-arg"
)

func main() {
	var cmds commands
	cmdLineArgs := os.Args[1:]

	path.SetupConfigDir()

	if len(cmdLineArgs) == 0 {
		fmt.Println("No arguments provided. Use -h or --help for more information.")
	} else {
		languages := templates.New()
		cmd := cmdLineArgs[0]

		if isReservedCommand(cmds, cmd) || cmd == "-h" || cmd == "--help" {
			arg.MustParse(&cmds)

			switch {
			case cmds.Create != nil:
				err := languages.CreateFileFromTemplate(cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)

				if err != nil {
					fmt.Println(err)
				}
			case cmds.List != nil:
				languages.List(cmds.List.Language)
			case cmds.Show != nil:
				err := languages.ShowTemplate(cmds.Show.Language, cmds.Show.TemplateName)

				if err != nil {
					fmt.Println(err)
				}
			case cmds.Run != nil:
				scriptToRun := scripts.New(languages)
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
