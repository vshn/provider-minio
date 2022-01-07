package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-logr/logr"
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
	// envPrefix is the global prefix to use for the keys in environment variables
	envPrefix = "BOOTSTRAP"
)

func main() {
	ctx, stop, app := newApp()
	defer stop()
	_ = app.RunContext(ctx, os.Args)
}

func newApp() (context.Context, context.CancelFunc, *cli.App) {
	logInstance := &atomic.Value{}
	logInstance.Store(logr.Discard())
	app := &cli.App{
		Name:     appName,
		Usage:    appLongName,
		Version:  fmt.Sprintf("%s, revision=%s, date=%s", version, commit, date),
		Compiled: compilationDate(),

		EnableBashCompletion: true,

		Before: beforeAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"verbose", "d"},
				Usage:   "sets the log level to debug",
				EnvVars: envVars("DEBUG"),
			},
			&cli.StringFlag{
				Name:        "log-format",
				Usage:       "sets the log format (values: [json, console])",
				EnvVars:     envVars("LOG_FORMAT"),
				DefaultText: "console",
			},
		},
		Commands: []*cli.Command{
			newExampleCommand(),
		},
		ExitErrHandler: func(context *cli.Context, err error) {
			if err != nil {
				AppLogger(context).Error(err, "fatal error")
				cli.HandleExitCoder(cli.Exit("", 1))
			}
		},
	}
	hasSubcommands := len(app.Commands) > 0
	app.Action = rootAction(hasSubcommands)
	// There is logr.NewContext(...) which returns a context that carries the logger instance.
	// However, since we are configuring and replacing this logger after starting up and parsing the flags,
	// we'll store a thread-safe atomic reference.
	parentCtx := context.WithValue(context.Background(), loggerContextKey{}, logInstance)
	ctx, stop := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM)
	return ctx, stop, app
}

func rootAction(hasSubcommands bool) func(context *cli.Context) error {
	return func(context *cli.Context) error {
		if hasSubcommands {
			return cli.ShowAppHelp(context)
		}
		logMetadata(context)
		return nil
	}
}

func beforeAction(c *cli.Context) error {
	setupLogging(c)
	if c.Args().Present() {
		// only print metadata if not displaying usage
		logMetadata(c)
	}
	return nil
}

func logMetadata(c *cli.Context) {
	log := AppLogger(c)
	if !usesProductionLoggingConfig(c) {
		log = log.WithValues("version", version)
	}
	log.WithValues(
		"date", date,
		"commit", commit,
		"go_os", runtime.GOOS,
		"go_arch", runtime.GOARCH,
		"go_version", runtime.Version(),
		"uid", os.Getuid(),
		"gid", os.Getgid(),
	).Info("Starting up " + appName)
}

// env combines envPrefix with given suffix delimited by underscore.
func env(suffix string) string {
	return envPrefix + "_" + suffix
}

// envVars combines envPrefix with each given suffix delimited by underscore.
func envVars(suffixes ...string) []string {
	arr := make([]string, len(suffixes))
	for i := range suffixes {
		arr[i] = env(suffixes[i])
	}
	return arr
}

func compilationDate() time.Time {
	compiled, err := time.Parse(time.RFC3339, date)
	if err != nil {
		// an empty Time{} causes cli.App to guess it from binary's file timestamp.
		return time.Time{}
	}
	return compiled
}
