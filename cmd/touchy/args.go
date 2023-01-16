package main

type ShowCommand struct {
	Language     string `arg:"positional,required" help:"The language that contains the template to show"`
	TemplateName string `arg:"positional" default:"default" help:"The name of the template to show"`
}

type ListCommand struct {
	Language string `arg:"positional" default:"all" help:"Get a list of templates for the given language"`
}

type CreateCommand struct {
	Language     string `arg:"positional,required" help:"language to use for template"`
	TemplateName string `arg:"positional" default:"default" help:"Name of the template to use"`
	FileName     string `arg:"positional" default:"DefaultOutputFileName" help:"Name of the generated file. Uses the key DefaultFileName in the language config file."`
}

type args struct {
	//TemplateName string       `arg:"positional required"`
	Create *CreateCommand `arg:"subcommand:create" help:"create a new template."`
	List   *ListCommand   `arg:"subcommand:list" help:"Show a list of all installed templates."`
	Show   *ShowCommand   `arg:"subcommand:show" help:"Show the contents of the template file."`
}

func (args) Description() string {
	return "Creates a file based upon a template"
}
