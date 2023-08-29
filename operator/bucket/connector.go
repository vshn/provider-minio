package bucket

import (
	"context"
	"fmt"
	"net/url"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/minio-go/v7"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	"github.com/vshn/provider-minio/operator/minioutil"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ managed.ExternalConnecter = &connector{}
var _ managed.ExternalClient = &bucketClient{}

const lockAnnotation = miniov1.Group + "/lock"

var (
	errNotBucket = fmt.Errorf("managed resource is not a bucket")
)

type connector struct {
	kube     client.Client
	recorder event.Recorder
}

type bucketClient struct {
	mc       *minio.Client
	recorder event.Recorder
}

// Connect implements managed.ExternalConnecter.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("connecting resource")

	bucket, ok := mg.(*miniov1.Bucket)
	if !ok {
		return nil, errNotBucket
	}

	config, err := c.getProviderConfig(ctx, bucket)
	if err != nil {
		return nil, err
	}

	mc, err := minioutil.NewMinioClient(ctx, c.kube, config)
	if err != nil {
		return nil, err
	}

	bc := &bucketClient{
		mc:       mc,
		recorder: c.recorder,
	}

	parsed, err := url.Parse(config.Spec.MinioURL)
	if err != nil {
		return nil, err
	}
	bucket.Status.Endpoint = parsed.Host
	bucket.Status.EndpointURL = parsed.String()

	return bc, nil
}

func (c *connector) getProviderConfig(ctx context.Context, bucket *miniov1.Bucket) (*providerv1.ProviderConfig, error) {
	configName := bucket.GetProviderConfigReference().Name
	config := &providerv1.ProviderConfig{}
	err := c.kube.Get(ctx, client.ObjectKey{Name: configName}, config)
	return config, err
}
