package user

import (
	"context"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (u *userClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {

	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("deleting resource")

	user, ok := mg.(*miniov1.User)
	if !ok {
		return managed.ExternalDelete{}, errNotUser
	}

	err := u.ma.RemoveUser(ctx, user.GetUserName())
	if err != nil {
		return managed.ExternalDelete{}, err
	}

	u.emitDeletionEvent(user)
	user.SetConditions(xpv1.Deleting())
	return managed.ExternalDelete{}, nil
}

func (u *userClient) emitDeletionEvent(user *miniov1.User) {
	u.recorder.Event(user, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Deleted",
		Message: "User successfully deleted",
	})
}
