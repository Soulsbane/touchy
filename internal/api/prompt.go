package api

import (
	"github.com/Iilun/survey/v2"
	"github.com/charmbracelet/huh"
)

type Prompts struct {
}

func NewPrompts() *Prompts {
	return &Prompts{}
}

func (p *Prompts) ConfirmationPrompt(message string) bool {
	name := false
	input := &survey.Confirm{
		Message: message,
		Default: false,
	}

	err := survey.AskOne(input, &name)

	if err != nil {
		return false
	}

	return name
}

func (p *Prompts) InputPrompt(message string, defaultValue string) string {
	value := ""
	input := &survey.Input{
		Message: message,
		Default: defaultValue,
	}

	err := survey.AskOne(input, &value)

	if err != nil {
		return defaultValue
	}

	return value

}

func (p *Prompts) MultiLineInputPrompt(message string, defaultValue string) string {
	value := defaultValue
	text := huh.NewText().
		Title(message).
		Value(&value)

	// FIXME: ctrl+enter does not submit the input
	// huh.NewForm(
	// 	huh.NewGroup(text),
	// ).WithKeyMap(&huh.KeyMap{
	// 	Text: huh.TextKeyMap{
	// 		NewLine: key.NewBinding(key.WithKeys("enter", "ctrl+j"), key.WithHelp("enter / ctrl+j", "new line")),
	// 		Editor:  key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "open editor")),
	// 		Submit:  key.NewBinding(key.WithKeys("ctrl+enter"), key.WithHelp("ctrl+enter", "submit")),
	// 	},
	// }).Run()
	huh.NewForm(huh.NewGroup(text)).Run()

	return value
}

func (p *Prompts) ChoicePrompt(message string, choices []string, defaultValue string) string {
	choice := ""
	input := &survey.Select{
		Message: message,
		Options: choices,
		Default: defaultValue,
	}

	err := survey.AskOne(input, &choice)

	if err != nil {
		return defaultValue
	}

	return choice
}

func (p *Prompts) MultiSelectPrompt(message string, promptChoices []string, defaultChoices []string) []string {
	choices := []string{}
	input := &survey.MultiSelect{
		Message: message,
		Options: promptChoices,
		Default: defaultChoices,
	}

	err := survey.AskOne(input, &choices)

	if err != nil {
		return defaultChoices
	}

	return choices
}
