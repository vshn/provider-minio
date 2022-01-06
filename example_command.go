package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type exampleCommand struct {
	ExampleFlag string
}

var exampleCommandName = "example"

func newExampleCommand() *cli.Command {
	command := &exampleCommand{}
	// TODO: Start hacking here
	return &cli.Command{
		Name:   exampleCommandName,
		Usage:  "Start example command",
		Before: command.validate,
		Action: command.execute,
		Flags: []cli.Flag{
			&cli.StringFlag{Destination: &command.ExampleFlag, Name: "flag", EnvVars: envVars("EXAMPLE_FLAG"), Value: "foo", Usage: "a demonstration how to configure the command"},
		},
	}
}

func (c *exampleCommand) validate(context *cli.Context) error {
	log := AppLogger(context).WithName(exampleCommandName)
	log.V(1).Info("validating config")
	if c.ExampleFlag == "" {
		return fmt.Errorf("option cannot be empty: %s", "flag")
	}
	return nil
}

func (c *exampleCommand) execute(context *cli.Context) error {
	log := AppLogger(context).WithName(exampleCommandName)
	log.Info("Hello from example command!", "config", c)
	return nil
}
