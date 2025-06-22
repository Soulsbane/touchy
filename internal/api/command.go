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

func (c *Command) RunWithOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (c *Command) RunWithOutputAsBytes(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	return output, nil
}
