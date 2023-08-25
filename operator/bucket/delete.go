package bucket

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (d *bucketClient) Delete(ctx context.Context, mg resource.Managed) error {
	log := controllerruntime.LoggerFrom(ctx)
	log.Info("deleting resource")

	return nil
}
