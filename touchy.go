package main

import (
	"github.com/Soulsbane/touchy/generator"
	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		//FileName     string `arg:"positional, required"`
		TemplateName string `arg:"positional, required"`
		FileName     string `arg:"-n, --name" default:"template" help:"Name of the generated file."`
	}

	arg.MustParse(&args)
	generator.CreateFileFromTemplate(args.FileName, args.TemplateName)
}
