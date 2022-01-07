package main

import (
	"fmt"
	"os"

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
		After:  command.shutdown,
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
	go func() {
		// This part enables graceful shutdowns. Can be removed if not needed.
		<-context.Done()
		err := c.shutdown(context)
		if err != nil {
			log.Error(err, "cannot properly shut down")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	log.Info("Hello from example command!", "config", c)
	return nil
}

func (c *exampleCommand) shutdown(context *cli.Context) error {
	log := AppLogger(context).WithName(exampleCommandName)
	log.Info("Shutting down example command")
	return nil
}
