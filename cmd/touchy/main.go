package main

import (
	"github.com/Soulsbane/touchy/internal/languages"
	"github.com/alexflint/go-arg"
)

func main() {
	var args args

	arg.MustParse(&args)
	languages := languages.New()

	switch {
	case args.Create != nil:
		languages.CreateFileFromTemplate(args.Create.FileName, args.Create.Language)
	case args.List != nil:
		languages.ListTemplates(args.List.Language)
	}

}
