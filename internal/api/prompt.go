package api

import (
	"errors"
	"fmt"
	"github.com/Iilun/survey/v2"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"os"
)

type Prompts struct {
}

func checkForPromptError(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Println("Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
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

func (p *Prompts) ChoicePromptWithDefault(message string, choices []string, defaultIndex int) string {
	choice, err := prompt.New().Ask(message).Choose(choices, choose.WithDefaultIndex(defaultIndex))
	checkForPromptError(err)

	return choice
}
