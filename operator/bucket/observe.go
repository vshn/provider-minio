package bucket

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"github.com/pkg/errors"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var bucketExistsFn = func(ctx context.Context, mc *minio.Client, bucketName string) (bool, error) {
	return mc.BucketExists(ctx, bucketName)
}

var bucketPolicyLatestFn = func(ctx context.Context, mc *minio.Client, bucketName string, policy string) (bool, error) {
	current, err := mc.GetBucketPolicy(ctx, bucketName)
	if err != nil {
		return false, err
	}

	return current == policy, nil
}

var bucketLifecycleLatestFn = func(ctx context.Context, mc *minio.Client, bucketName string, lifecycleRules *lifecycle.Configuration) (bool, error) {
	current, err := mc.GetBucketLifecycle(ctx, bucketName)
	// Continuing if error is not "The lifecycle configuration does not exist" as we want to report the resource as not up-to-date, if that is the case
	if err != nil && !strings.Contains(err.Error(), "The lifecycle configuration does not exist") {
		return false, err
	}

	return current == lifecycleRules, nil
}

func (d *bucketClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	log := controllerruntime.LoggerFrom(ctx)
	log.V(1).Info("observing resource")

	bucket, ok := mg.(*miniov1.Bucket)
	if !ok {
		return managed.ExternalObservation{}, errNotBucket
	}

	bucketName := bucket.GetBucketName()
	exists, err := bucketExistsFn(ctx, d.mc, bucketName)

	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.StatusCode == http.StatusForbidden {
			// As we have full control over the minio instance, we can say with confidence that this case is a
			// "permission denied"
			return managed.ExternalObservation{}, errors.Wrap(err, "permission denied, please check the provider-config")
		}
		if errResp.StatusCode == http.StatusMovedPermanently {
			return managed.ExternalObservation{}, errors.Wrap(err, "mismatching endpointURL and zone, or bucket exists already in a different region, try changing bucket name")
		}
		return managed.ExternalObservation{}, errors.Wrap(err, "cannot determine whether bucket exists")
	}
	if _, hasAnnotation := bucket.GetAnnotations()[lockAnnotation]; hasAnnotation && exists {
		bucket.Status.AtProvider.BucketName = bucketName
		bucket.SetConditions(xpv1.Available())

		isLatest := true
		if bucket.Spec.ForProvider.Policy != nil {
			u, err := bucketPolicyLatestFn(ctx, d.mc, bucketName, *bucket.Spec.ForProvider.Policy)
			if err != nil {
				return managed.ExternalObservation{}, errors.Wrap(err, "cannot determine whether a bucket policy exists")
			}

			isLatest = u
		}

		if isLatest && bucket.Spec.ForProvider.LifecycleRules != nil {
			lifecycleConfiguration := lifecycle.NewConfiguration()
			for _, rule := range bucket.Spec.ForProvider.LifecycleRules {
				lifecycleRule := lifecycle.Rule{
					ID: rule.ID,
					Expiration: lifecycle.Expiration{
						Days: lifecycle.ExpirationDays(rule.ExpirationDays),
					},
					NoncurrentVersionExpiration: lifecycle.NoncurrentVersionExpiration{
						NoncurrentDays: lifecycle.ExpirationDays(rule.NoncurrentVersionExpirationDays),
					},
					Status: "Enabled",
				}
				lifecycleConfiguration.Rules = append(lifecycleConfiguration.Rules, lifecycleRule)

				upToDate, err := bucketLifecycleLatestFn(ctx, d.mc, bucketName, lifecycleConfiguration)
				if err != nil {
					return managed.ExternalObservation{}, errors.Wrap(err, "cannot determine whether a bucket lifecycle rule exists")
				}

				if !upToDate {
					isLatest = false
					break
				}
			}
		}
		return managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: isLatest}, nil
	} else if exists {
		return managed.ExternalObservation{}, fmt.Errorf("bucket already exists, try changing bucket name: %s", bucketName)
	}

	return managed.ExternalObservation{}, nil
}
