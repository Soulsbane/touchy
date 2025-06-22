package api

import "os/exec"

type Command struct {
}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
