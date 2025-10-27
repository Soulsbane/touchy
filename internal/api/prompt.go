package api

import (
	"github.com/Iilun/survey/v2"
	"github.com/charmbracelet/bubbles/key"
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
	input := huh.NewInput().Title(message).Value(&value)
	err := huh.NewForm(huh.NewGroup(input)).Run()

	if err != nil {
		return defaultValue
	}

	return value
}

func (p *Prompts) MultiLineInputPrompt(message string, defaultValue string) string {
	value := ""
	text := huh.NewText().Title(message).Value(&value)

	err := huh.NewForm(
		huh.NewGroup(text),
	).WithKeyMap(&huh.KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		Text: huh.TextKeyMap{
			NewLine: key.NewBinding(key.WithKeys("enter", "ctrl+j"), key.WithHelp("enter / ctrl+j", "new line")),
			// Editor:  key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "open editor")),
			Submit: key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "submit")),
		},
	}).Run()

	if err != nil {
		return defaultValue
	}

	return value
}

func (p *Prompts) ChoicePrompt(message string, choices []string, defaultValue string) string {
	choice := ""
	err := huh.NewSelect[string]().
		Options(huh.NewOptions(choices...)...).
		Value(&choice).
		Title(message).Run()

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
