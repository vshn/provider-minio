package user

import (
	"context"
	"strings"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/madmin-go/v3"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (u *userClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	log := controllerruntime.LoggerFrom(ctx)
	log.V(1).Info("updating resource")

	user, ok := mg.(*miniov1.User)
	if !ok {
		return managed.ExternalUpdate{}, errNotUser
	}

	userInfo, err := u.ma.GetUserInfo(ctx, user.GetUserName())
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	policies := strings.Split(userInfo.PolicyName, ",")
	// empty string will just result in an array with one empty string, so let's check that
	if policies[0] == "" {
		policies = []string{}
	}

	if len(policies) > 0 {
		req := madmin.PolicyAssociationReq{
			Policies: policies,
			User:     user.GetUserName(),
		}
		_, err = u.ma.DetachPolicy(ctx, req)
		if err != nil {
			return managed.ExternalUpdate{}, err
		}
	}

	err = u.setUserPolicies(ctx, user.GetUserName(), user.Spec.ForProvider.Policies)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	u.emitUpdateEvent(user)
	return managed.ExternalUpdate{}, nil
}

func (u *userClient) emitUpdateEvent(user *miniov1.User) {
	u.recorder.Event(user, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Updated",
		Message: "User successfully updated",
	})
}
