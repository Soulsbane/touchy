package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Soulsbane/touchy/internal/templates"
	"github.com/alexflint/go-arg"
)

func isReservedCommand(cmds commands, command string) bool {
	// This checks if a command is reserved based on the command existing as a member of the commands struct.
	dummyVal := reflect.ValueOf(cmds)
	numFields := dummyVal.NumField()

	for i := 0; i < numFields; i++ {
		field := dummyVal.Type().Field(i).Name

		if strings.ToLower(field) == command {
			return true
		}
	}

	return false
}

func main() {
	var cmds commands
	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) == 0 {
		fmt.Println("No arguments provided. Use -h or --help for more information.")
	} else {
		languages := templates.New()
		cmd := cmdLineArgs[0]

		if isReservedCommand(cmds, cmd) {
			arg.MustParse(&cmds)

			switch {
			case cmds.Create != nil:
				languages.CreateFileFromTemplate(cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)
			case cmds.List != nil:
				languages.List(cmds.List.Language)
			case cmds.Show != nil:
				languages.ShowTemplate(cmds.Show.Language, cmds.Show.TemplateName)
			}
		} else {
			var createCmd CreateCommand

			arg.MustParse(&createCmd)
			languages.CreateFileFromTemplate(createCmd.Language, createCmd.TemplateName, createCmd.FileName)
		}
	}
}
