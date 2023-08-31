package user

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/madmin-go/v3"
	"github.com/sethvargo/go-password/password"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	// UserCreatedAnnotationKey is the annotation name where we store the information that the
	// user has been created.
	UserCreatedAnnotationKey string = "minio.crossplane.io/user-created"
)

func (u *userClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {

	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("creating resource")

	user, ok := mg.(*miniov1.User)
	if !ok {
		return managed.ExternalCreation{}, errNotUser
	}

	secretKey, err := password.Generate(64, 5, 0, false, true)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	// The minioAdmin doesn't return an error if the user already exists, it just overrides it...
	exists, err := u.userExists(ctx, user.GetUserName())
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	if exists {
		return managed.ExternalCreation{}, fmt.Errorf("user already exists")
	}

	err = u.ma.AddUser(ctx, user.GetUserName(), secretKey)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	err = u.setUserPolicies(ctx, user.GetUserName(), user.Spec.ForProvider.Policies)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	u.emitCreationEvent(user)

	annotations := user.GetAnnotations()
	annotations[UserCreatedAnnotationKey] = "true"
	user.SetAnnotations(annotations)

	connectionDetails := managed.ConnectionDetails{
		SecretKeyName: []byte(secretKey),
		AccessKeyName: []byte(user.GetUserName()),
	}

	return managed.ExternalCreation{ConnectionDetails: connectionDetails}, nil
}

func (u *userClient) userExists(ctx context.Context, name string) (bool, error) {
	users, err := u.ma.ListUsers(ctx)
	if err != nil {
		return false, err
	}

	_, exists := users[name]
	return exists, nil
}

func (u *userClient) emitCreationEvent(user *miniov1.User) {
	u.recorder.Event(user, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Created",
		Message: "User successfully created",
	})
}

func (u *userClient) setUserPolicies(ctx context.Context, userName string, policies []string) error {

	if len(policies) == 0 {
		return nil
	}

	req := madmin.PolicyAssociationReq{
		Policies: policies,
		User:     userName,
	}

	_, err := u.ma.AttachPolicy(ctx, req)

	return err
}
