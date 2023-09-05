package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/pkg/bucket/policy"
	miniopolicy "github.com/minio/pkg/iam/policy"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type jsonPolicy []byte

const (
	// PolicyCreatedAnnotationKey is the annotation name where we store the information that the
	// user has been created.
	PolicyCreatedAnnotationKey string = "minio.crossplane.io/policy-created"
)

func (p *policyClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {

	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("creating resource")

	policy, ok := mg.(*miniov1.Policy)
	if !ok {
		return managed.ExternalCreation{}, errNotPolicy
	}

	policyies, err := p.ma.ListCannedPolicies(ctx)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	if _, ok := policyies[policy.GetName()]; ok {
		return managed.ExternalCreation{}, fmt.Errorf("policy already exists")
	}

	if policy.Spec.ForProvider.AllowBucket != "" {
		return managed.ExternalCreation{}, p.createBucketPolicy(ctx, policy)
	}

	if policy.Spec.ForProvider.RawPolicy != "" {
		return managed.ExternalCreation{}, p.createRawPolicy(ctx, policy)
	}

	return managed.ExternalCreation{}, fmt.Errorf("no policy specified")
}

func (p *policyClient) createBucketPolicy(ctx context.Context, policy *miniov1.Policy) error {
	parsedPolicy, err := p.getAllowBucketPolicy(policy.Spec.ForProvider.AllowBucket)
	if err != nil {
		return err
	}

	err = p.ma.AddCannedPolicy(ctx, policy.GetName(), parsedPolicy)
	if err != nil {
		return err
	}

	p.emitCreationEvent(policy)
	p.setLock(policy)

	return nil
}

func (p *policyClient) createRawPolicy(ctx context.Context, policy *miniov1.Policy) error {
	err := p.ma.AddCannedPolicy(ctx, policy.GetName(), []byte(policy.Spec.ForProvider.RawPolicy))
	if err != nil {
		return err
	}

	p.emitCreationEvent(policy)
	p.setLock(policy)

	return nil
}

func (p *policyClient) getAllowBucketPolicy(bucket string) (jsonPolicy, error) {

	actionSet := miniopolicy.NewActionSet(miniopolicy.AllActions)

	resourceSet := miniopolicy.NewResourceSet(
		miniopolicy.NewResource(bucket, "/"),
		miniopolicy.NewResource(bucket, "*"),
	)

	newPolicy := miniopolicy.Policy{
		Version: "2012-10-17",
		Statements: []miniopolicy.Statement{
			{
				SID:       "addPerm",
				Effect:    policy.Allow,
				Actions:   actionSet,
				Resources: resourceSet,
			},
		},
	}

	err := newPolicy.Validate()
	if err != nil {
		return nil, err
	}

	return json.Marshal(newPolicy)
}

func (p *policyClient) emitCreationEvent(policy *miniov1.Policy) {
	p.recorder.Event(policy, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Created",
		Message: "User successfully created",
	})
}

func (p *policyClient) setLock(policy *miniov1.Policy) {
	annotations := policy.GetAnnotations()
	annotations[PolicyCreatedAnnotationKey] = "claimed"
	policy.SetAnnotations(annotations)
}
