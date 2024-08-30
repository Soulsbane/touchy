package main

import (
	"fmt"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/templates"
	"github.com/Soulsbane/touchy/internal/ui"
	"github.com/jedib0t/go-pretty/v6/table"
	"reflect"
	"strings"
)

type ShowCommand struct {
	Language     string `arg:"positional,required" help:"The language that contains the template to show"`
	TemplateName string `arg:"positional" default:"default" help:"The name of the template to show"`
}

type ListCommand struct {
	Type     string `arg:"positional" default:"all" help:"Option to list either scripts, templates or all"`
	Language string `arg:"positional" default:"all" help:"Get a list of templates or scripts for the given language"`
}

type CreateCommand struct {
	Language     string `arg:"positional,required" help:"language to use for template"`
	TemplateName string `arg:"positional" default:"default" help:"name of the template to use"`
	FileName     string `arg:"positional" default:"DefaultOutputFileName" help:"name of the generated file. Uses the key DefaultFileName in the language config file."`
}

type RunCommand struct {
	ScriptName string `arg:"positional,required" help:"name of the script to run"`
}

type commands struct {
	// TemplateName string       `arg:"positional required"`
	Create *CreateCommand `arg:"subcommand:create" help:"create a new template."`
	List   *ListCommand   `arg:"subcommand:list" help:"Show a list of all installed templates."`
	Show   *ShowCommand   `arg:"subcommand:show" help:"Show the contents of the template file."`
	Run    *RunCommand    `arg:"subcommand:run" help:"Run a script."`
}

func (commands) Description() string {
	return "Creates a file based upon a template"
}
func isReservedCommand(cmds commands, command string) bool {
	// This checks if a command is reserved based on the command existing as a member of the command's struct.
	dummyVal := reflect.ValueOf(cmds)
	numFields := dummyVal.NumField()

	for i := 0; i < numFields; i++ {
		field := dummyVal.Type().Field(i).Name

		if strings.ToLower(field) == command {
			return true
		}
	}

	return false
}

func ListLanguages() {
	temps := templates.New()
	languages := temps.GetListOfAllLanguages()

	outputTable := ui.CreateNewTableWriter("Languages", "Name", "Long Name", "Description")

	for name, language := range languages {
		languageInfo := language.GetInfoFile()

		outputTable.AppendRow(table.Row{name, languageInfo.GetName(), languageInfo.GetDescription()})
	}

	outputTable.Render()
}

func ListScripts(scripts []infofile.InfoFile) {
	outputTable := ui.CreateNewTableWriter("Scripts", "Script name", "Description")

	for _, script := range scripts {
		outputTable.AppendRow(table.Row{script.GetName(), script.GetDescription()})
	}

	outputTable.Render()
}

func ListTemplates(listArg string) {
	temps := templates.New()
	languages := temps.GetListOfAllLanguages()

	if languageTemplates, found := languages[listArg]; found {
		languageInfo := languageTemplates.GetInfoFile()
		outputTable := ui.CreateNewTableWriter(languageInfo.GetName()+" Templates", "name", "Description", "Default Output File name")

		for _, info := range languageTemplates.GetTemplatesInfoFiles() {
			outputTable.AppendRow(table.Row{info.GetName(), info.GetDescription(), info.GetDefaultOutputFileName()})
		}

		outputTable.Render()
	} else if listArg == "all" {
		for _, language := range languages {
			languageInfo := language.GetInfoFile()
			outputTable := ui.CreateNewTableWriter(languageInfo.GetName()+" Templates", "name", "Description", "Default Output File name")

			for _, info := range language.GetTemplatesInfoFiles() {
				outputTable.AppendRow(table.Row{info.GetName(), info.GetDescription(), info.GetDefaultOutputFileName()})
			}

			outputTable.Render()
		}
	} else {
		fmt.Println("That language could not be found! Use 'list all' to see all available languages.")
	}
}
