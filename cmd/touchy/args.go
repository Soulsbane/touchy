package main

import (
	"github.com/Soulsbane/touchy/internal/infofile"
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
	TemplateName string `arg:"positional" default:"default" help:"Name of the template to use"`
	FileName     string `arg:"positional" default:"DefaultOutputFileName" help:"Name of the generated file. Uses the key DefaultFileName in the language config file."`
}

type RunCommand struct {
	ScriptName string `arg:"positional,required" help:"Name of the script to run"`
}

type commands struct {
	//TemplateName string       `arg:"positional required"`
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

func ListScripts(scripts []infofile.InfoFile) {
	outputTable := ui.CreateNewTableWriter("Scripts", "Script Name", "Description")

	for _, script := range scripts {
		outputTable.AppendRow(table.Row{script.Name, script.Description})
	}

	outputTable.Render()
}
