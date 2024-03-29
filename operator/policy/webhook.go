package policy

import (
	"context"
	"fmt"

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

	policy, ok := obj.(*miniov1.Policy)
	if !ok {
		return nil, errNotPolicy
	}

	return nil, v.validatePolicy(policy)
}

// ValidateUpdate implements admission.CustomValidator.
func (v *Validator) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	v.log.V(1).Info("Validate update")

	newPolicy, ok := newObj.(*miniov1.Policy)
	if !ok {
		return nil, errNotPolicy
	}

	return nil, v.validatePolicy(newPolicy)
}

// ValidateDelete implements admission.CustomValidator.
func (v *Validator) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	v.log.V(1).Info("validate delete (noop)")
	return nil, nil
}

func (v *Validator) validatePolicy(policy *miniov1.Policy) error {
	if policy.Spec.ForProvider.AllowBucket != "" && policy.Spec.ForProvider.RawPolicy != "" {
		return fmt.Errorf(".spec.forProvider.allowBucket and .spec.forProvider.rawPolicy are mutual exclusive, please only specify one")
	}

	providerConfigRef := policy.Spec.ProviderConfigReference
	if providerConfigRef == nil || providerConfigRef.Name == "" {
		return field.Invalid(field.NewPath("spec", "providerConfigRef", "name"), "null", "Provider config is required")
	}
	return nil
}
