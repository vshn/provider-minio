package bucket

import (
	"context"
	"net/http"
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/go-logr/logr"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestProvisioningPipeline_Observe(t *testing.T) {
	policy := "policy-struct"
	tests := map[string]struct {
		givenBucket  *miniov1.Bucket
		bucketExists bool
		returnError  error
		policyLatest bool

		expectedError             string
		expectedResult            managed.ExternalObservation
		expectedBucketObservation miniov1.BucketProviderStatus
	}{
		"NewBucketDoesntYetExistOnMinio": {
			givenBucket: &miniov1.Bucket{Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
				BucketName: "my-bucket"}},
			},
			expectedResult: managed.ExternalObservation{},
		},
		"NewBucketWithPolicyDoesntYetExistOnMinio": {
			givenBucket: &miniov1.Bucket{Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
				BucketName: "my-bucket-with-policy",
				Policy:     &policy}},
			},
			expectedResult: managed.ExternalObservation{},
		},
		"BucketExistsAndAccessibleWithOurCredentials": {
			givenBucket: &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{
					lockAnnotation: "claimed",
				}},
				Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
					BucketName: "my-bucket"}},
			},
			bucketExists:              true,
			expectedResult:            managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true},
			expectedBucketObservation: miniov1.BucketProviderStatus{BucketName: "my-bucket"},
		},
		"NewBucketObservationThrowsGenericError": {
			givenBucket: &miniov1.Bucket{Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
				BucketName: "my-bucket"}},
			},
			returnError:    errors.New("error"),
			expectedResult: managed.ExternalObservation{},
			expectedError:  "cannot determine whether bucket exists: error",
		},
		"BucketAlreadyExistsOnMinio_WithoutAccess": {
			givenBucket: &miniov1.Bucket{Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
				BucketName: "my-bucket"}},
			},
			returnError:    minio.ErrorResponse{StatusCode: http.StatusForbidden, Message: "Access Denied"},
			expectedResult: managed.ExternalObservation{},
			expectedError:  "permission denied, please check the provider-config: Access Denied",
		},
		"BucketAlreadyExistsOnMinio_WithAccess_PreventAdoption": {
			// this is a case where we should avoid adopting an existing bucket even if we have access.
			// Otherwise, there could be multiple K8s resources that manage the same bucket.
			givenBucket: &miniov1.Bucket{
				Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
					BucketName: "my-bucket"}},
				// no bucket name in status here.
			},
			bucketExists:   true,
			expectedResult: managed.ExternalObservation{},
			expectedError:  "bucket already exists, try changing bucket name: my-bucket",
		},
		"BucketAlreadyExistsOnMinio_InAnotherZone": {
			givenBucket: &miniov1.Bucket{
				Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
					BucketName: "my-bucket"}},
			},
			returnError:    minio.ErrorResponse{StatusCode: http.StatusMovedPermanently, Message: "301 Moved Permanently"},
			expectedResult: managed.ExternalObservation{},
			expectedError:  "mismatching endpointURL and zone, or bucket exists already in a different region, try changing bucket name: 301 Moved Permanently",
		},
		"BucketPolicyNoChangeRequired": {
			givenBucket: &miniov1.Bucket{
				ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{
					lockAnnotation: "claimed",
				}},
				Spec: miniov1.BucketSpec{ForProvider: miniov1.BucketParameters{
					BucketName: "my-bucket",
					Policy:     &policy}},
			},
			policyLatest:              true,
			bucketExists:              true,
			expectedResult:            managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true},
			expectedBucketObservation: miniov1.BucketProviderStatus{BucketName: "my-bucket"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			currFn := bucketExistsFn
			defer func() {
				bucketExistsFn = currFn
			}()

			bucketPolicyLatestFn = func(ctx context.Context, mc *minio.Client, bucketName string, policy string) (bool, error) {
				return tc.policyLatest, tc.returnError
			}

			bucketExistsFn = func(ctx context.Context, mc *minio.Client, bucketName string) (bool, error) {
				return tc.bucketExists, tc.returnError
			}
			b := bucketClient{}
			result, err := b.Observe(logr.NewContext(context.Background(), logr.Discard()), tc.givenBucket)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedBucketObservation, tc.givenBucket.Status.AtProvider)
		})
	}
}
