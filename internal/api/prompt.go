package api

import (
	"github.com/Iilun/survey/v2"
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
