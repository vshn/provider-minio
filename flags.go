package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	// Remove `-v` short option from --version flag in favor of verbosity.
	cli.VersionFlag.(*cli.BoolFlag).Aliases = nil
}

func newLogLevelFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name: "log-level", Aliases: []string{"v"}, EnvVars: envVars("LOG_LEVEL"),
		Usage: "number of the log level verbosity",
		Value: 0,
	}
}

func newLogFormatFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name: "log-format", EnvVars: envVars("LOG_FORMAT"),
		Usage: "sets the log format (values: [json, console])",
		Value: "console",
		Action: func(context *cli.Context, format string) error {
			if format == "console" || format == "json" {
				return nil
			}
			_ = cli.ShowAppHelp(context)
			return fmt.Errorf("unknown log format: %s", format)
		},
	}
}

func newExampleFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name: "flag", EnvVars: envVars("EXAMPLE_FLAG"), Required: true,
		Usage:       "a demonstration how to configure the command",
		Destination: dest,
		Action: func(context *cli.Context, s string) error {
			if len(s) >= 3 {
				return nil
			}
			_ = cli.ShowAppHelp(context)
			return fmt.Errorf("option needs at least 3 characters: %s", "flag")
		},
	}
}
