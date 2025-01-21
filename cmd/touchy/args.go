package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/ui"
	"github.com/jedib0t/go-pretty/v6/table"
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

func handleCreateCommand(languageName string, templateName string, fileName string) {
	err := manager.CreateFileFromTemplate(languageName, templateName, fileName)

	if err != nil {
		handleError(err, templateName, languageName)
	}
}

func handleListCommand(listType string, languageName string) {
	switch listType {
	case "all":
		scriptsList := touchyScripts.GetListOfScripts()
		ListScripts(scriptsList)
		fmt.Println("")
		manager.ListTemplates(languageName)
	case "languages":
		manager.ListLanguages()
	case "scripts":
		scriptsList := touchyScripts.GetListOfScripts()
		ListScripts(scriptsList)
	case "templates":
		manager.ListTemplates(languageName)
	default:
		// NOTE: listType in this case could be a language name
		if manager.HasLanguage(listType) {
			manager.ListTemplates(listType)
		} else {
			// TODO: Add support for pulling the default language  template
			fmt.Println("That list type could not be found! Use 'list all' to see all available types.")
		}
	}
}

func handleShowCommand(languageName string, templateName string) {
	err := manager.ShowTemplate(languageName, templateName)

	if err != nil {
		handleError(err, templateName, languageName)
	}
}

func handleRunCommand(scriptName string) {
	err := touchyScripts.Run(scriptName)

	if err != nil {
		handleError(err, scriptName, "")
	}
}

func ListScripts(scripts []infofile.InfoFile) {
	outputTable := ui.CreateNewTableWriter("Scripts", "Script name", "Description")

	for _, script := range scripts {
		outputTable.AppendRow(table.Row{script.GetName(), script.GetDescription()})
	}

	outputTable.Render()
}
