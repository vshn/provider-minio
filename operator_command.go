package main

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
	"github.com/vshn/provider-minio/apis"
	"github.com/vshn/provider-minio/operator"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type operatorCommand struct {
	LeaderElectionEnabled bool
	WebhookCertDir        string
}

func newOperatorCommand() *cli.Command {
	command := &operatorCommand{}
	return &cli.Command{
		Name:   "operator",
		Usage:  "Start provider operator in mode",
		Action: command.execute,
		Flags: []cli.Flag{
			newLeaderElectionEnabledFlag(&command.LeaderElectionEnabled),
			newWebhookTLSCertDirFlag(&command.WebhookCertDir),
		},
	}
}

func (o *operatorCommand) execute(ctx *cli.Context) error {
	_ = LogMetadata(ctx)
	log := logr.FromContextOrDiscard(ctx.Context).WithName(ctx.Command.Name)
	log.Info("Starting up operator", "config", o)
	ctrl.SetLogger(log)

	cfg, err := ctrl.GetConfig()
	if err != nil {
		return err
	}

	// configure client-side throttling
	cfg.QPS = 100
	cfg.Burst = 150 // more Openshift friendly

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		// controller-runtime uses both ConfigMaps and Leases for leader election by default.
		// Leases expire after 15 seconds, with a 10-second renewal deadline.
		// We've observed leader loss due to renewal deadlines being exceeded when under high load - i.e.
		//  hundreds of reconciles per second and ~200rps to the API server.
		// Switching to Leases only and longer leases appears to alleviate this.
		LeaderElection:             o.LeaderElectionEnabled,
		LeaderElectionID:           "leader-election-provider-minio",
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaseDuration:              func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:              func() *time.Duration { d := 50 * time.Second; return &d }(),
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    9443,
			CertDir: o.WebhookCertDir,
		}),
	})

	if err != nil {
		return err
	}

	err = apis.AddToScheme(mgr.GetScheme())
	if err != nil {
		return err
	}

	if o.WebhookCertDir != "" {
		err = operator.SetupWebhooks(mgr)
		if err != nil {
			return err
		}
	}

	err = operator.SetupControllers(mgr)
	if err != nil {
		return err
	}

	return mgr.Start(ctx.Context)
}
