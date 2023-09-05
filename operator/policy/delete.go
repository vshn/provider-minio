package policy

import (
	"context"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	metav1 "github.com/vshn/provider-minio/apis/minio/v1"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (p *policyClient) Delete(ctx context.Context, mg resource.Managed) error {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("deleting resource")

	policy, ok := mg.(*metav1.Policy)
	if !ok {
		return errNotPolicy
	}

	policy.SetConditions(xpv1.Deleting())
	p.emitDeletionEvent(policy)
	return p.ma.RemoveCannedPolicy(ctx, policy.GetName())
}

func (p *policyClient) emitDeletionEvent(policy *miniov1.Policy) {
	p.recorder.Event(policy, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Deleted",
		Message: "Policy successfully deleted",
	})
}
