package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
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
	cli.VersionPrinter = func(_ *cli.Context) {
		fmt.Printf("version=%s revision=%s date=%s\n", version, commit, date)
	}

	compiled, err := time.Parse(time.RFC3339, date)
	if err != nil {
		compiled = time.Time{}
	}

	app := &cli.App{
		Name:     appName,
		Usage:    appLongName,
		Version:  version,
		Compiled: compiled,

		EnableBashCompletion: true,

		Before: before,
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
	app.Action = func(context *cli.Context) error {
		if hasSubcommands {
			return cli.ShowAppHelp(context)
		}
		logMetadata(AppLogger(context))
		return nil
	}
	ctx, stop := signal.NotifyContext(context.WithValue(context.Background(), AppContextKeyName, app), syscall.SIGINT, syscall.SIGTERM)
	return ctx, stop, app
}

func before(c *cli.Context) error {
	useProductionConfig := strings.EqualFold("JSON", c.String("log-format"))
	logger := newLogger(appName, c.Bool("debug"), useProductionConfig)
	SetAppLogger(c, logger)
	if !c.Args().Present() {
		// skip printing metadata if displaying the usage
		return nil
	}
	log := logger.WithValues()
	if !useProductionConfig {
		log = log.WithValues("version", version)
	}
	logMetadata(logger)
	return nil
}

func logMetadata(log logr.Logger) {
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
