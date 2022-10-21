package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	// these variables are populated by Goreleaser when releasing
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")

	// TODO: Adjust app name
	appName     = "go-bootstrap"
	appLongName = "a generic bootstrapping project"

	// TODO: Adjust or clear env var prefix
	// envPrefix is the global prefix to use for the keys in environment variables.
	// Include a delimiter like `_` if required.
	envPrefix = "BOOTSTRAP_"
)

func main() {
	app := newApp()
	err := app.Run(os.Args)
	// If required flags aren't set, it will return with error before we could set up logging
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := &cli.App{
		Name:    appName,
		Usage:   appLongName,
		Version: fmt.Sprintf("%s, revision=%s, date=%s", version, commit, date),

		Before: setupLogging,
		Flags: []cli.Flag{
			newLogLevelFlag(),
			newLogFormatFlag(),
		},
		Commands: []*cli.Command{
			newExampleCommand(),
		},
	}
	return app
}

// TODO: Remove env() and envVars() if not using an environment variable prefix.

// env combines envPrefix with given suffix delimited by underscore.
func env(suffix string) string {
	return envPrefix + suffix
}

// envVars combines envPrefix with each given suffix delimited by underscore.
func envVars(suffixes ...string) []string {
	arr := make([]string, len(suffixes))
	for i := range suffixes {
		arr[i] = env(suffixes[i])
	}
	return arr
}
