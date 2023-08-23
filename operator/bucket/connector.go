package bucket

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ managed.ExternalConnecter = &connector{}
var _ managed.ExternalClient = &bucketClient{}

type connector struct {
	kube         client.Client
	recorder     event.Recorder
	bucketClient *bucketClient
}

type bucketClient struct {
}

// Connect implements managed.ExternalConnecter.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("connecting resource")

	c.bucketClient = &bucketClient{}

	return c.bucketClient, nil
}
