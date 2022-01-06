package main

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

const (
	signalHandlerMetadataKeyName  = "signalHandler"
	currentContextMetadataKeyName = "context"
)

func terminate(app *cli.App) {
	log := app.Metadata[loggerMetadataKeyName].(logr.Logger)
	if h, found := app.Metadata[signalHandlerMetadataKeyName]; found {
		err := h.(func(context *cli.Context) error)(app.Metadata[currentContextMetadataKeyName].(*cli.Context))
		if err != nil {
			log.Error(err, "cannot properly shut down app")
			os.Exit(1)
		}
	}
	log.Info("Shutting down")
	os.Exit(0)
}

// SetSignalHandler sets the given handler.
// The handler will be invoked with given context in case of SIGINT or SIGTERM.
// This setter should be called in an Action function (Before or Action itself) of a cli.App.
func SetSignalHandler(ctx *cli.Context, handler func(context *cli.Context) error) {
	ctx.App.Metadata[signalHandlerMetadataKeyName] = handler
	ctx.App.Metadata[currentContextMetadataKeyName] = ctx
}
