package policy

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (p *policyClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {

	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("update resource")

	policy, ok := mg.(*miniov1.Policy)
	if !ok {
		return managed.ExternalUpdate{}, errNotPolicy
	}

	policies, err := p.ma.ListCannedPolicies(ctx)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	_, ok = policies[policy.GetName()]
	if !ok {
		return managed.ExternalUpdate{}, fmt.Errorf("policy does not exist")
	}

	if policy.Spec.ForProvider.AllowBucket != "" {
		p.emitUpdateEvent(policy)
		return managed.ExternalUpdate{}, p.createBucketPolicy(ctx, policy)
	}

	if policy.Spec.ForProvider.RawPolicy != "" {
		p.emitUpdateEvent(policy)
		return managed.ExternalUpdate{}, p.createRawPolicy(ctx, policy)
	}

	return managed.ExternalUpdate{}, nil
}

func (p *policyClient) emitUpdateEvent(policy *miniov1.Policy) {
	p.recorder.Event(policy, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Updated",
		Message: "Policy successfully updated",
	})
}
