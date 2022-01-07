package main

import (
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// LoggerMetadataKeyName is the key which can be used to retrieve the logr.Logger from cli.App's Metadata.
	// Usage example:
	//  app.Metadata[LoggerMetadataKeyName].(logr.Logger)
	LoggerMetadataKeyName = "cli:logger"
	// AppContextKeyName is the key which can be used to retrieve the cli.App from context.Context.
	// Usage example:
	//	ctx.Value(AppContextKeyName).(*cli.App)
	AppContextKeyName = "cli:app"
)

// AppLogger retrieves the application-wide logger instance from the cli.Context's Metadata.
// This function will return nil if SetAppLogger was not called before this function is called.
func AppLogger(c *cli.Context) logr.Logger {
	return c.App.Metadata[LoggerMetadataKeyName].(logr.Logger)
}

// SetAppLogger stores the application-wide logger instance to the cli.Context's Metadata,
// so that it can later be retrieved by AppLogger.
func SetAppLogger(c *cli.Context, logger logr.Logger) {
	c.App.Metadata[LoggerMetadataKeyName] = logger
}

func newLogger(name string, debug bool, useProductionConfig bool) logr.Logger {
	zc := zap.NewDevelopmentConfig()
	zc.EncoderConfig.ConsoleSeparator = " | "
	if useProductionConfig {
		zc = zap.NewProductionConfig()
	}
	if debug {
		// Zap's levels get more verbose as the number gets smaller,
		// bug logr's level increases with greater numbers.
		zc.Level = zap.NewAtomicLevelAt(zapcore.Level(-2)) // max logger.V(2)
	}
	z, err := zc.Build()
	zap.ReplaceGlobals(z)
	if err != nil {
		log.Fatalf("error configuring the logging stack")
	}
	logger := zapr.NewLogger(z).WithName(name)
	if useProductionConfig {
		// Append the version to each log so that logging stacks like EFK/Loki can correlate errors with specific versions.
		return logger.WithValues("version", version)
	}
	return logger
}
