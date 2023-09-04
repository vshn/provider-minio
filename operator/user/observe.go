package user

import (
	"context"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/madmin-go/v3"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
)

const (
	AccessKeyName = "AWS_ACCESS_KEY_ID"
	SecretKeyName = "AWS_SECRET_ACCESS_KEY"
)

func (u *userClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {

	user, ok := mg.(*miniov1.User)
	if !ok {
		return managed.ExternalObservation{}, errNotUser
	}

	_, ok = user.GetAnnotations()[UserCreatedAnnotationKey]
	if !ok && user.Status.AtProvider.UserName == "" {
		// The user has not yet been create, let's do it then
		return managed.ExternalObservation{}, nil
	}

	user.Status.AtProvider.UserName = user.GetUserName()

	users, err := u.ma.ListUsers(ctx)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	minioUser, ok := users[user.GetUserName()]
	if !ok {
		// The user doesn't exist!
		// Let's try again.
		user.Status.AtProvider.UserName = ""
		return managed.ExternalObservation{}, nil
	}

	user.Status.AtProvider.Status = string(minioUser.Status)

	if minioUser.Status == madmin.AccountEnabled {
		user.SetConditions(xpv1.Available())
	}

	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true}, nil
}
