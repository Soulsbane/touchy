package main

import (
	"github.com/Soulsbane/touchy/internal/generator"
	"github.com/alexflint/go-arg"
)

func main() {
	var args args

	arg.MustParse(&args)
	generator := generator.New()

	switch {
	case args.Create != nil:
		generator.CreateFileFromTemplate(args.Create.FileName, args.Create.Language)
	case args.List != nil:
		generator.ListTemplates()
	}

}
