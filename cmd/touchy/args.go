package main

import "github.com/logrusorgru/aurora"

type ShowCommand struct {
}

type ListCommand struct {
}

type CreateCommand struct {
	Language string `arg:"positional" help:"language to use for template"`
}

type args struct {
	//TemplateName string       `arg:"positional required"`
	Create   *CreateCommand `arg:"subcommand:create" help:"create a new template"`
	FileName string         `arg:"-n, --name" default:"template" help:"Name of the generated file. If not specified, the filename will be the config's DefaultFileName value."`
	List     *ListCommand   `arg:"subcommand:list" help:"Show a list of all installed templates"`
	Show     *ShowCommand   `arg:"subcommand:show" help:"Show the contents of the template file"`
}

func (args) Description() string {
	return aurora.Blue("Creates a file based upon a template.").Bold().String()
}
