package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	"github.com/vshn/provider-minio/operator/minioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	_                   admission.CustomValidator = &Validator{}
	getProviderConfigFn                           = getProviderConfig
	getMinioAdminFn                               = getMinioAdmin
)

type cannedPolicyLister interface {
	ListCannedPolicies(context.Context) (map[string]json.RawMessage, error)
}

// Validator validates admission requests.
type Validator struct {
	log  logr.Logger
	kube client.Client
}

// ValidateCreate implements admission.CustomValidator.
func (v *Validator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	v.log.V(1).Info("Validate create")

	user, ok := obj.(*miniov1.User)
	if !ok {
		return nil, errNotUser
	}

	providerConfigRef := user.Spec.ProviderConfigReference
	if providerConfigRef == nil || providerConfigRef.Name == "" {
		return nil, field.Invalid(field.NewPath("spec", "providerConfigRef", "name"), "null", "Provider config is required")
	}

	err := v.doesPolicyExist(ctx, user)
	if err != nil {
		return nil, field.Invalid(field.NewPath("spec", "forProvider", "policies"), user.Spec.ForProvider.Policies, err.Error())
	}

	return nil, nil
}

// ValidateUpdate implements admission.CustomValidator.
func (v *Validator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	v.log.V(1).Info("Validate update")

	oldUser, ok := oldObj.(*miniov1.User)
	if !ok {
		return nil, errNotUser
	}
	newUser, ok := newObj.(*miniov1.User)
	if !ok {
		return nil, errNotUser
	}

	if newUser.GetUserName() != oldUser.GetUserName() {
		return nil, field.Invalid(field.NewPath("spec", "forProvider", "userName"), newUser.GetUserName(), "Changing the username is not allowed")
	}

	providerConfigRef := newUser.Spec.ProviderConfigReference
	if providerConfigRef == nil || providerConfigRef.Name == "" {
		return nil, field.Invalid(field.NewPath("spec", "providerConfigRef", "name"), "null", "Provider config is required")
	}

	// Deleting an object also triggers an update first to set the
	// deletionTimestamp. If the policy and the user objects are deleted at the
	// same time, it's possible that the policy has already been gone by the time
	// the user update happens.
	// So we ignore checking if the timestamp is set.
	if newUser.GetDeletionTimestamp() != nil {
		return nil, nil
	}

	err := v.doesPolicyExist(ctx, newUser)
	if err != nil {
		return nil, field.Invalid(field.NewPath("spec", "forProvider", "policies"), newUser.Spec.ForProvider.Policies, err.Error())
	}

	return nil, nil
}

// ValidateDelete implements admission.CustomValidator.
func (v *Validator) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	v.log.V(1).Info("validate delete (noop)")
	return nil, nil
}

func (v *Validator) doesPolicyExist(ctx context.Context, user *miniov1.User) error {

	if len(user.Spec.ForProvider.Policies) == 0 {
		return nil
	}

	config, err := getProviderConfigFn(ctx, user, v.kube)
	if err != nil {
		return err
	}

	ma, err := getMinioAdminFn(ctx, v.kube, config)
	if err != nil {
		return err
	}

	policies, err := ma.ListCannedPolicies(ctx)
	if err != nil {
		return err
	}

	for _, policy := range user.Spec.ForProvider.Policies {
		_, ok := policies[policy]
		if !ok {
			return fmt.Errorf("policy not found: %s", policy)
		}
	}

	return nil
}

func getProviderConfig(ctx context.Context, user *miniov1.User, kube client.Client) (*providerv1.ProviderConfig, error) {
	configName := user.GetProviderConfigReference().Name
	config := &providerv1.ProviderConfig{}
	err := kube.Get(ctx, client.ObjectKey{Name: configName}, config)
	return config, err
}

func getMinioAdmin(ctx context.Context, kube client.Client, config *providerv1.ProviderConfig) (cannedPolicyLister, error) {
	return minioutil.NewMinioAdmin(ctx, kube, config)
}
