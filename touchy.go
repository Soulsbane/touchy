package main

import (
	"github.com/alexflint/go-arg"
)

func loadTemplate(templateName string) {
}

func main() {
	var args struct {
		TemplateName string `arg:"positional, required"`
	}

	arg.MustParse(&args)
	loadTemplate(args.TemplateName)
}
