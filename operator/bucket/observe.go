package bucket

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (d *bucketClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	log := controllerruntime.LoggerFrom(ctx)
	log.V(1).Info("observing resource")

	return managed.ExternalObservation{}, nil
}
