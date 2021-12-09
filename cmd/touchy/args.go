package main

import "github.com/logrusorgru/aurora"

type args struct {
	TemplateName string `arg:"positional, required"`
	FileName     string `arg:"-n, --name" help:"Name of the generated file. If not specified, the filename will be the config's DefaultFileName value."`
	//Show         string `arg:"subcommand:show" help:"Show the contents of the template file"`
	//List         string `arg:"subcommand:list" help:"Show a list of all installed templates"`
}

func (args) Description() string {
	return aurora.Blue("Creates a file based upon a template.").Bold().String()
}
