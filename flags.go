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
		Name: "log-level", Aliases: []string{"v"}, EnvVars: []string{"LOG_LEVEL"},
		Usage: "number of the log level verbosity",
		Value: 0,
	}
}

func newLogFormatFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name: "log-format", EnvVars: []string{"LOG_FORMAT"},
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

func newLeaderElectionEnabledFlag(dest *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name: "leader-election-enabled", Value: false, EnvVars: []string{"LEADER_ELECTION_ENABLED"},
		Usage:       "Use leader election for the controller manager.",
		Destination: dest,
	}
}

func newWebhookTLSCertDirFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name: "webhook-tls-cert-dir", EnvVars: []string{"WEBHOOK_TLS_CERT_DIR"}, // Env var is set by Crossplane
		Usage:       "Directory containing the certificates for the webhook server. If empty, the webhook server is not started.",
		Destination: dest,
	}
}
