package example

import (
	"github.com/urfave/cli/v2"
	"github.com/vshn/go-bootstrap/cfg"
	"github.com/vshn/go-bootstrap/cmd"
)

var (
	// TODO: Start hacking here

	// Command configures the CLI subcommand
	Command = &cli.Command{
		Name:   "example",
		Usage:  "Start example command",
		Action: main,
		Flags: []cli.Flag{
			&cli.StringFlag{Destination: &cfg.Config.ExampleFlag, Name: "flag", EnvVars: []string{cfg.Env("EXAMPLE_FLAG")}, Value: "foo", Usage: "an demonstration how to configure the subcommand"},
		},
	}
)

func main(c *cli.Context) error {
	log := cmd.AppLogger(c).WithName("example")
	log.Info("Hello from example command!", "config", cfg.Config)
	return nil
}
