package main

import (
	"testing"

	"github.com/go-logr/zapr"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap/zaptest"
)

func newTestApp(t *testing.T) *cli.App {
	return &cli.App{
		Metadata: map[string]interface{}{
			LoggerMetadataKeyName: zapr.NewLogger(zaptest.NewLogger(t)),
		},
	}
}

func newAppContext(t *testing.T) *cli.Context {
	return cli.NewContext(newTestApp(t), nil, nil)
}
