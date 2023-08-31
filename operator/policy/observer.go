package policy

import (
	"context"
	"encoding/json"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	miniopolicy "github.com/minio/pkg/iam/policy"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (p *policyClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {

	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("observing resource")

	policy, ok := mg.(*miniov1.Policy)
	if !ok {
		return managed.ExternalObservation{}, errNotPolicy
	}

	_, ok = policy.GetAnnotations()[PolicyCreatedAnnotationKey]
	if !ok {
		// The policy has not yet been create, let's do it then
		return managed.ExternalObservation{}, nil
	}

	policies, err := p.ma.ListCannedPolicies(ctx)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	observedPolicy, ok := policies[policy.GetName()]
	if !ok {
		// The policy hasn't yet been created it seems
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	if policy.Spec.ForProvider.AllowBucket != "" {
		bucketPolicy, err := p.getAllowBucketPolicy(policy.Spec.ForProvider.AllowBucket)
		if err != nil {
			return managed.ExternalObservation{}, err
		}

		equal, err := p.sameObject(json.RawMessage(bucketPolicy), observedPolicy)
		if err != nil {
			return managed.ExternalObservation{}, err
		}
		if !equal {
			policy.SetConditions(miniov1.Updating())
			return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: false}, nil
		}
	}

	if policy.Spec.ForProvider.RawPolicy != "" {
		equal, err := p.sameObject(json.RawMessage(policy.Spec.ForProvider.RawPolicy), observedPolicy)
		if err != nil {
			return managed.ExternalObservation{}, err
		}
		if !equal {
			policy.SetConditions(miniov1.Updating())
			return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: false}, nil
		}
	}

	policy.Status.AtProvider.Policy = string(observedPolicy)
	policy.SetConditions(xpv1.Available())

	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true}, nil
}

// sameObject will marshal both given objects to a map.
// After that it will do a deepEquals to verify that they have the equal values.
func (p *policyClient) sameObject(a, b json.RawMessage) (bool, error) {
	aPol := &miniopolicy.Policy{}
	bPol := &miniopolicy.Policy{}

	err := json.Unmarshal(a, aPol)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(b, bPol)
	if err != nil {
		return false, err
	}

	return aPol.Equals(*bPol), nil
}
