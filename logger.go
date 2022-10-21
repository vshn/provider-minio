package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogMetadata prints various metadata to the root logger.
// It prints version, architecture and current user ID and returns nil.
func LogMetadata(c *cli.Context) error {
	log := logr.FromContextOrDiscard(c.Context)
	log.WithValues(
		"version", version,
		"date", date,
		"commit", commit,
		"go_os", runtime.GOOS,
		"go_arch", runtime.GOARCH,
		"go_version", runtime.Version(),
		"uid", os.Getuid(),
		"gid", os.Getgid(),
	).Info("Starting up " + appName)
	return nil
}

func setupLogging(c *cli.Context) error {
	log, err := newZapLogger(appName, c.Int(newLogLevelFlag().Name), usesProductionLoggingConfig(c))
	c.Context = logr.NewContext(c.Context, log)
	return err
}

func usesProductionLoggingConfig(c *cli.Context) bool {
	return strings.EqualFold("JSON", c.String(newLogFormatFlag().Name))
}

func newZapLogger(name string, verbosityLevel int, useProductionConfig bool) (logr.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.ConsoleSeparator = " | "
	if useProductionConfig {
		cfg = zap.NewProductionConfig()
	}
	// Zap's levels get more verbose as the number gets smaller,
	// bug logr's level increases with greater numbers.
	cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(verbosityLevel * -1))
	z, err := cfg.Build()
	if err != nil {
		return logr.Discard(), fmt.Errorf("error configuring the logging stack: %w", err)
	}
	zap.ReplaceGlobals(z)
	zlog := zapr.NewLogger(z).WithName(name)
	if useProductionConfig {
		// Append the version to each log so that logging stacks like EFK/Loki can correlate errors with specific versions.
		return zlog.WithValues("version", version), nil
	}
	return zlog, nil
}
