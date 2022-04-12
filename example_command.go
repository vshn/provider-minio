package main

import (
	"fmt"
	"os"
	"sync"

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
			&cli.StringFlag{Destination: &command.ExampleFlag, Name: "flag", EnvVars: envVars("EXAMPLE_FLAG"), Value: "foo", Usage: "a demonstration how to configure the command", Required: true},
		},
	}
}

func (c *exampleCommand) validate(ctx *cli.Context) error {
	_ = LogMetadata(ctx)
	log := AppLogger(ctx).WithName(exampleCommandName)
	log.V(1).Info("validating config")
	// The `Required` property in the StringFlag above already checks if it's non-empty.
	if len(c.ExampleFlag) <= 2 {
		return fmt.Errorf("option needs at least 3 characters: %s", "flag")
	}
	return nil
}

func (c *exampleCommand) execute(ctx *cli.Context) error {
	log := AppLogger(ctx).WithName(exampleCommandName)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// Shutdown hook. Can be used to gracefully shutdown listeners or pre-shutdown cleanup.
		// Can be removed if not needed.
		// Please note that this example is incomplete and doesn't cover all cases when properly implementing shutdowns.
		defer wg.Done()
		<-ctx.Done()
		err := c.shutdown(ctx)
		if err != nil {
			log.Error(err, "cannot properly shut down")
			os.Exit(2)
		}
	}()
	log.Info("Hello from example command!", "config", c)
	wg.Wait()
	return nil
}

func (c *exampleCommand) shutdown(ctx *cli.Context) error {
	log := AppLogger(ctx).WithName(exampleCommandName)
	log.Info(fmt.Sprintf("Shutting down %q command", exampleCommandName))
	return nil
}
