package api

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

type Prompts struct {
}

func NewPrompts() *Prompts {
	return &Prompts{}
}

func (p *Prompts) ConfirmationPrompt(message string) bool {
	value := false
	err := huh.NewConfirm().Title(message).Affirmative("Yes").Negative("No").Value(&value).Run()

	if err != nil {
		return false
	}

	return value
}
func (p *Prompts) InputPrompt(message string, defaultValue string) string {
	value := ""
	err := huh.NewInput().Title(message).Value(&value).Run()

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
			Submit:  key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "submit")),
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
	err := huh.NewMultiSelect[string]().
		Options(huh.NewOptions(promptChoices...)...).
		Value(&choices).
		Title(message).Run()

	if err != nil {
		return defaultChoices
	}

	return choices
}
