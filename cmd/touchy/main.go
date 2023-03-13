package main

import (
	"fmt"
	"os"

	"github.com/Soulsbane/touchy/internal/templates"
	"github.com/alexflint/go-arg"
)

func main() {
	var args args
	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) == 0 {
		fmt.Println("No arguments provided. Use -h or --help for more information.")
	}

	arg.MustParse(&args)
	languages := templates.New()

	switch {
	case args.Create != nil:
		languages.CreateFileFromTemplate(args.Create.Language, args.Create.TemplateName, args.Create.FileName)
	case args.List != nil:
		languages.List(args.List.Language)
	case args.Show != nil:
		languages.ShowTemplate(args.Show.Language, args.Show.TemplateName)
	}

}
