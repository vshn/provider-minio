package policy

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/madmin-go/v3"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	"github.com/vshn/provider-minio/operator/minioutil"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	errNotPolicy = fmt.Errorf("managed resource is not a policy")
)

type connector struct {
	kube     client.Client
	recorder event.Recorder
}

type policyClient struct {
	ma       *madmin.AdminClient
	recorder event.Recorder
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("connecting resource")

	policy, ok := mg.(*miniov1.Policy)
	if !ok {
		return nil, errNotPolicy
	}

	config, err := c.getProviderConfig(ctx, policy)
	if err != nil {
		return nil, err
	}

	ma, err := minioutil.NewMinioAdmin(ctx, c.kube, config)
	if err != nil {
		return nil, err
	}

	uc := &policyClient{
		ma:       ma,
		recorder: c.recorder,
	}

	return uc, nil
}

func (c *connector) getProviderConfig(ctx context.Context, Policy *miniov1.Policy) (*providerv1.ProviderConfig, error) {
	configName := Policy.GetProviderConfigReference().Name
	config := &providerv1.ProviderConfig{}
	err := c.kube.Get(ctx, client.ObjectKey{Name: configName}, config)
	return config, err
}
