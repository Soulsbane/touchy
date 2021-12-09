package main

import (
	"github.com/Soulsbane/touchy/internal/generator"
	"github.com/alexflint/go-arg"
)

func main() {
	var args args

	arg.MustParse(&args)

	switch {
	case args.Create != nil:
		generator.CreateFileFromTemplate(args.FileName, args.Create.Language)
	case args.List != nil:
		generator.ListTemplates()
	}

}
