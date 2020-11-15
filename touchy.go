package main

import (
	"github.com/Soulsbane/touchy/generator"
	"github.com/alexflint/go-arg"
	"github.com/logrusorgru/aurora"
)

type args struct {
	TemplateName string `arg:"positional, required"`
	FileName     string `arg:"-n, --name" default:"template" help:"Name of the generated file."`
}

func (args) Description() string {
	return aurora.Blue("Creates a file based upon a template.").Bold().String()
}

func main() {
	var args args

	arg.MustParse(&args)
	generator.CreateFileFromTemplate(args.FileName, args.TemplateName)
}
