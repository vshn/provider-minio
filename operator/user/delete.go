package user

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (u *userClient) Delete(ctx context.Context, mg resource.Managed) error {

	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("deleting resource")

	user, ok := mg.(*miniov1.User)
	if !ok {
		return errNotUser
	}

	err := u.ma.RemoveUser(ctx, user.GetUserName())
	if err != nil {
		return err
	}

	u.emitDeletionEvent(user)
	return nil
}

func (u *userClient) emitDeletionEvent(user *miniov1.User) {
	u.recorder.Event(user, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Deleted",
		Message: "User successfully deleted",
	})
}
