package user

import (
	"context"

	"github.com/go-logr/logr"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.CustomValidator = &Validator{}

// Validator validates admission requests.
type Validator struct {
	log logr.Logger
}

// ValidateCreate implements admission.CustomValidator.
func (v *Validator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	v.log.V(1).Info("Validate create")

	user, ok := obj.(*miniov1.User)
	if !ok {
		return nil, errNotUser
	}

	providerConfigRef := user.Spec.ProviderConfigReference
	if providerConfigRef == nil || providerConfigRef.Name == "" {
		return nil, field.Invalid(field.NewPath("spec", "providerConfigRef", "name"), "null", "Provider config is required")
	}

	return nil, nil
}

// ValidateUpdate implements admission.CustomValidator.
func (v *Validator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
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

	return nil, nil
}

// ValidateDelete implements admission.CustomValidator.
func (v *Validator) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	v.log.V(1).Info("validate delete (noop)")
	return nil, nil
}
