package main

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap/zaptest"
)

func newAppContext(t *testing.T) *cli.Context {
	logger := zapr.NewLogger(zaptest.NewLogger(t))
	instance := &atomic.Value{}
	instance.Store(logger)
	return cli.NewContext(&cli.App{}, nil, &cli.Context{
		Context: context.WithValue(context.Background(), loggerContextKey{}, instance),
	})
}
