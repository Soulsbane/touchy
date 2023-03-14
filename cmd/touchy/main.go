package main

import (
	"fmt"
	"os"

	"github.com/Soulsbane/touchy/internal/templates"
	"github.com/alexflint/go-arg"
)

func main() {
	var cmds commands
	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) == 0 {
		fmt.Println("No arguments provided. Use -h or --help for more information.")
	}

	arg.MustParse(&cmds)
	languages := templates.New()

	switch {
	case cmds.Create != nil:
		languages.CreateFileFromTemplate(cmds.Create.Language, cmds.Create.TemplateName, cmds.Create.FileName)
	case cmds.List != nil:
		languages.List(cmds.List.Language)
	case cmds.Show != nil:
		languages.ShowTemplate(cmds.Show.Language, cmds.Show.TemplateName)
	}

}
