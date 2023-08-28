package bucket

import (
	"context"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidator_ValidateCreate_RequireProviderConfig(t *testing.T) {
	tests := map[string]struct {
		providerName  string
		expectedError string
	}{
		"GivenProviderName_ThenExpectNoError": {
			providerName: "provider-config",
		},
		"GivenNoProviderName_ThenExpectError": {
			providerName:  "",
			expectedError: `.spec.providerConfigRef.name is required`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bucket := &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: miniov1.BucketSpec{
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{Name: tc.providerName},
					},
					ForProvider: miniov1.BucketParameters{BucketName: "bucket"},
				},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateCreate(context.TODO(), bucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateUpdate_PreventBucketNameChange(t *testing.T) {
	tests := map[string]struct {
		newBucketName string
		oldBucketName string
		expectedError string
	}{
		"GivenNoNameInStatus_WhenNoNameInSpec_ThenExpectNil": {
			oldBucketName: "",
			newBucketName: "",
		},
		"GivenNoNameInStatus_WhenNameInSpec_ThenExpectNil": {
			oldBucketName: "",
			newBucketName: "my-bucket",
		},
		"GivenNameInStatus_WhenNameInSpecSame_ThenExpectNil": {
			oldBucketName: "my-bucket",
			newBucketName: "my-bucket",
		},
		"GivenNameInStatus_WhenNameInSpecEmpty_ThenExpectNil": {
			oldBucketName: "bucket",
			newBucketName: "", // defaults to metadata.name
		},
		"GivenNameInStatus_WhenNameInSpecDifferent_ThenExpectError": {
			oldBucketName: "my-bucket",
			newBucketName: "different",
			expectedError: `spec.forProvider.bucketName: Invalid value: "different": Changing the bucket name is not allowed after creation`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldBucket := &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: miniov1.BucketSpec{
					ForProvider:  miniov1.BucketParameters{BucketName: tc.oldBucketName},
					ResourceSpec: xpv1.ResourceSpec{ProviderConfigReference: &xpv1.Reference{Name: "provider-config"}},
				},
				Status: miniov1.BucketStatus{AtProvider: miniov1.BucketProviderStatus{BucketName: tc.oldBucketName}},
			}
			newBucket := &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: miniov1.BucketSpec{
					ForProvider:  miniov1.BucketParameters{BucketName: tc.newBucketName},
					ResourceSpec: xpv1.ResourceSpec{ProviderConfigReference: &xpv1.Reference{Name: "provider-config"}},
				},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateUpdate(context.TODO(), oldBucket, newBucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateUpdate_RequireProviderConfig(t *testing.T) {
	tests := map[string]struct {
		providerConfigToRef *xpv1.Reference
		expectedError       string
	}{
		"GivenProviderConfigRefWithName_ThenExpectNoError": {
			providerConfigToRef: &xpv1.Reference{
				Name: "provider-config",
			},
		},
		"GivenProviderConfigEmptyRef_ThenExpectError": {
			providerConfigToRef: &xpv1.Reference{
				Name: "",
			},
			expectedError: `spec.providerConfigRef.name: Invalid value: "null": Provider config is required`,
		},
		"GivenProviderConfigRefNil_ThenExpectError": {
			providerConfigToRef: nil,
			expectedError:       `spec.providerConfigRef.name: Invalid value: "null": Provider config is required`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldBucket := &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: miniov1.BucketSpec{
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: tc.providerConfigToRef,
					},
					ForProvider: miniov1.BucketParameters{BucketName: "bucket"},
				},
				Status: miniov1.BucketStatus{AtProvider: miniov1.BucketProviderStatus{BucketName: "bucket"}},
			}
			newBucket := &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: miniov1.BucketSpec{
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: tc.providerConfigToRef,
					},
					ForProvider: miniov1.BucketParameters{BucketName: "bucket"},
				},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateUpdate(context.TODO(), oldBucket, newBucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateUpdate_PreventZoneChange(t *testing.T) {
	tests := map[string]struct {
		newZone       string
		oldZone       string
		expectedError string
	}{
		"GivenZoneUnchanged_ThenExpectNil": {
			oldZone: "zone",
			newZone: "zone",
		},
		"GivenZoneChanged_ThenExpectError": {
			oldZone:       "zone",
			newZone:       "different",
			expectedError: `spec.forProvider.region: Invalid value: "different": Changing the region is not allowed after creation`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldBucket := &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: miniov1.BucketSpec{
					ForProvider:  miniov1.BucketParameters{Region: tc.oldZone},
					ResourceSpec: xpv1.ResourceSpec{ProviderConfigReference: &xpv1.Reference{Name: "provider-config"}},
				},
				Status: miniov1.BucketStatus{AtProvider: miniov1.BucketProviderStatus{BucketName: "bucket"}},
			}
			newBucket := &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Name: "bucket"},
				Spec: miniov1.BucketSpec{
					ForProvider:  miniov1.BucketParameters{Region: tc.newZone},
					ResourceSpec: xpv1.ResourceSpec{ProviderConfigReference: &xpv1.Reference{Name: "provider-config"}},
				},
				Status: miniov1.BucketStatus{AtProvider: miniov1.BucketProviderStatus{BucketName: "bucket"}},
			}
			v := &Validator{log: logr.Discard()}
			_, err := v.ValidateUpdate(context.TODO(), oldBucket, newBucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
