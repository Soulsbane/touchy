package api

import (
	"errors"
	"fmt"
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
	choice, err := prompt.New().Ask(message).Choose([]string{"Yes", "No"}, choose.WithDefaultIndex(1))
	checkForPromptError(err)

	if choice == "No" {
		return false
	} else {
		return true
	}
}

func (p *Prompts) InputPrompt(message string, defaultValue string) string {
	value, err := prompt.New().Ask(message).Input(defaultValue)
	checkForPromptError(err)

	return value
}

func (p *Prompts) ChoicePrompt(message string, choices []string) string {
	choice, err := prompt.New().Ask(message).Choose(choices)
	checkForPromptError(err)

	return choice
}

func (p *Prompts) ChoicePromptWithDefault(message string, choices []string, defaultIndex int) string {
	choice, err := prompt.New().Ask(message).Choose(choices, choose.WithDefaultIndex(defaultIndex))
	checkForPromptError(err)

	return choice
}
