package main

import (
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type exampleCommand struct {
	ExampleFlag string
}

func newExampleCommand() *cli.Command {
	command := &exampleCommand{}
	// TODO: Start hacking here
	return &cli.Command{
		Name:   "example",
		Usage:  "Start example command",
		Before: command.validate,
		Action: command.execute,
		Flags: []cli.Flag{
			newExampleFlag(&command.ExampleFlag),
		},
	}
}

func (c *exampleCommand) validate(ctx *cli.Context) error {
	log := logr.FromContextOrDiscard(ctx.Context).WithName(ctx.Command.Name)
	log.V(1).Info("validating config")
	return nil
}

func (c *exampleCommand) execute(ctx *cli.Context) error {
	_ = LogMetadata(ctx)
	log := logr.FromContextOrDiscard(ctx.Context).WithName(ctx.Command.Name)
	log.Info("Hello from example command!", "config", c)
	return nil
}
