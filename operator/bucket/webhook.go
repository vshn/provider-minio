package bucket

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
	bucket, ok := obj.(*miniov1.Bucket)
	if !ok {
		return nil, errNotBucket
	}
	v.log.V(1).Info("Validate create", "name", bucket.Name)

	providerConfigRef := bucket.Spec.ProviderConfigReference
	if providerConfigRef == nil || providerConfigRef.Name == "" {
		return nil, fmt.Errorf(".spec.providerConfigRef.name is required")
	}
	return nil, nil
}

// ValidateUpdate implements admission.CustomValidator.
func (v *Validator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	newBucket := newObj.(*miniov1.Bucket)
	oldBucket := oldObj.(*miniov1.Bucket)
	v.log.V(1).Info("Validate update")

	if oldBucket.Status.AtProvider.BucketName != "" {
		if newBucket.GetBucketName() != oldBucket.Status.AtProvider.BucketName {
			return nil, field.Invalid(field.NewPath("spec", "forProvider", "bucketName"), newBucket.Spec.ForProvider.BucketName, "Changing the bucket name is not allowed after creation")
		}
		if newBucket.Spec.ForProvider.Region != oldBucket.Spec.ForProvider.Region {
			return nil, field.Invalid(field.NewPath("spec", "forProvider", "region"), newBucket.Spec.ForProvider.Region, "Changing the region is not allowed after creation")
		}
	}
	providerConfigRef := newBucket.Spec.ProviderConfigReference
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
