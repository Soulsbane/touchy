package main

import (
	"github.com/Soulsbane/touchy/internal/generator"
	"github.com/alexflint/go-arg"
)

func main() {
	var args args

	arg.MustParse(&args)
	generator.CreateFileFromTemplate(args.FileName, args.TemplateName)
}
