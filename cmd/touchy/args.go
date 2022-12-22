package main

import "github.com/logrusorgru/aurora"

type ShowCommand struct {
}

type ListCommand struct {
	Language string `arg:"positional" help:"Get a list of templates for the given language"`
}

type CreateCommand struct {
	Language string `arg:"positional,required" help:"language to use for template"`
	FileName string `arg:"positional" help:"Name of the generated file. Uses the key DefaultFileName in the language config file."`
}

type args struct {
	//TemplateName string       `arg:"positional required"`
	Create *CreateCommand `arg:"subcommand:create" help:"create a new template."`
	List   *ListCommand   `arg:"subcommand:list" help:"Show a list of all installed templates."`
	Show   *ShowCommand   `arg:"subcommand:show" help:"Show the contents of the template file."`
}

func (args) Description() string {
	return aurora.Blue("Creates a file based upon a template.").Bold().String()
}
