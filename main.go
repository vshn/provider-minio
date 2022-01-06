package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/urfave/cli/v2"
	"github.com/vshn/go-bootstrap/cfg"
	"github.com/vshn/go-bootstrap/cmd"
	"github.com/vshn/go-bootstrap/cmd/example"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// these variables are populated by Goreleaser when releasing
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")

	// TODO: Adjust app name
	appName     = "go-bootstrap"
	appLongName = "a generic bootstrapping project"
)

func main() {
	err := app().Run(os.Args)
	if err != nil {
		log.Fatalf("unable to start %s: %v", appName, err)
	}
}

func before(c *cli.Context) error {
	logger := newLogger(appName, c.Bool("debug"))
	cmd.SetAppLogger(c, logger)
	if c.Args().Present() {
		logger.WithValues(
			"version", version,
			"date", date,
			"commit", commit,
			"go_os", runtime.GOOS,
			"go_arch", runtime.GOARCH,
			"go_version", runtime.Version(),
			"uid", os.Getuid(),
			"gid", os.Getgid(),
		).Info("Starting up " + appName)
	}

	return nil
}

func app() *cli.App {
	cli.VersionPrinter = func(_ *cli.Context) {
		fmt.Printf("version=%s revision=%s date=%s\n", version, commit, date)
	}

	compiled, err := time.Parse(time.RFC3339, date)
	if err != nil {
		compiled = time.Time{}
	}

	return &cli.App{
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
				EnvVars: []string{cfg.Env("DEBUG")},
			},
		},
		Commands: []*cli.Command{
			example.Command,
		},
	}
}

func newLogger(name string, debug bool) logr.Logger {
	zc := zap.NewDevelopmentConfig()
	if debug {
		// Zap's levels get more verbose as the number gets smaller,
		// bug logr's level increases with greater numbers.
		zc.Level = zap.NewAtomicLevelAt(zapcore.Level(-2)) // max logger.V(2)
	}
	z, err := zc.Build(zap.WithCaller(false))
	if err != nil {
		log.Fatalf("error configuring the logging stack")
	}
	return zapr.NewLogger(z).WithName(name)
}
