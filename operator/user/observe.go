package user

import (
	"context"
	"reflect"
	"strings"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	k8svi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	AccessKeyName = "AWS_ACCESS_KEY_ID"
	SecretKeyName = "AWS_SECRET_ACCESS_KEY"
)

func (u *userClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	log := ctrl.LoggerFrom(ctx)

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
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	if !u.equalPolicies(minioUser, user) {
		user.SetConditions(miniov1.Updating())
		return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: false}, nil
	}

	user.Status.AtProvider.Status = string(minioUser.Status)
	user.Status.AtProvider.Policies = minioUser.PolicyName

	if minioUser.Status == madmin.AccountEnabled {
		user.SetConditions(xpv1.Available())
	} else {
		user.SetConditions(miniov1.Disabled())
	}

	if mg.GetDeletionTimestamp() == nil {

		secret := k8svi.Secret{}

		err = u.kube.Get(ctx, types.NamespacedName{
			Namespace: mg.GetWriteConnectionSecretToReference().Namespace,
			Name:      mg.GetWriteConnectionSecretToReference().Name,
		}, &secret)
		if err != nil {
			return managed.ExternalObservation{}, err
		}

		mclient, err := minio.New(u.url.Host, &minio.Options{
			Creds:  credentials.NewStaticV4(string(secret.Data[AccessKeyName]), string(secret.Data[SecretKeyName]), ""),
			Secure: u.tlsSettings,
		})
		if err != nil {
			return managed.ExternalObservation{ResourceUpToDate: false, ResourceExists: true}, nil
		}

		_, err = mclient.ListBuckets(context.Background())
		// AccessDenied is ok in this context, because we just want to check if the user has working credentials
		if err != nil && err.Error() != "Access Denied." {
			return managed.ExternalObservation{ResourceUpToDate: false, ResourceExists: true}, nil
		}

		log.Info("user client created, everything went fine " + string(secret.Data[AccessKeyName]) + " " + string(secret.Data[SecretKeyName]))
	}

	return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true}, nil
}

func (u *userClient) equalPolicies(minioUser madmin.UserInfo, user *miniov1.User) bool {
	// policyName contains a string with all applied policies seperated by comma
	minioPolicies := strings.Split(minioUser.PolicyName, ",")

	// if policyName is an empty string, then string.Split() will create an array with just one empty string
	// to make it comparable we need to catch this case and set it to an empty array
	if minioPolicies[0] == "" {
		minioPolicies = nil
	}

	return reflect.DeepEqual(minioPolicies, user.Spec.ForProvider.Policies)
}
