package bucket

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var encryptionTypeMap map[string]string = map[string]string{
	"sse-kms": "aws:kms",
	"sse-s3":  "AES256",
}

func (b *bucketClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	log := controllerruntime.LoggerFrom(ctx)
	log.V(1).Info("updating resource")

	bucket, ok := mg.(*miniov1.Bucket)
	if !ok {
		return managed.ExternalUpdate{}, errNotBucket
	}

	if bucket.Spec.ForProvider.Policy != nil {
		err := b.mc.SetBucketPolicy(ctx, bucket.GetBucketName(), *bucket.Spec.ForProvider.Policy)
		if err != nil {
			return managed.ExternalUpdate{}, err
		}
	}

	// Encryption
	bucketName := bucket.GetBucketName()
	if do, sseConfig, err := createSseConfig(bucket); err != nil {
		return managed.ExternalUpdate{}, err
	} else if do {
		currentConfig, err := b.mc.GetBucketEncryption(ctx, bucketName)
		if err != nil {
			return managed.ExternalUpdate{}, fmt.Errorf("could not get bucket encryption config: %w", err)
		}

		// Only update encryption if it actually changed
		thisEncType := encryptionTypeMap[bucket.Spec.ForProvider.Encryption.Type]
		if currentConfig.Rules[0].Apply.SSEAlgorithm != thisEncType {
			if err := b.mc.SetBucketEncryption(ctx, bucketName, &sseConfig); err != nil {
				return managed.ExternalUpdate{}, fmt.Errorf("failed updating encryption: %w", err)
			}
		}
	} else {
		if err := b.mc.RemoveBucketEncryption(ctx, bucketName); err != nil {
			return managed.ExternalUpdate{}, fmt.Errorf("failed removing encryption: %w", err)
		}
	}

	return managed.ExternalUpdate{}, nil
}
