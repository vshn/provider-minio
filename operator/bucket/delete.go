package bucket

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/minio-go/v7"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (b *bucketClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	log := controllerruntime.LoggerFrom(ctx)
	log.Info("deleting resource")

	bucket, ok := mg.(*miniov1.Bucket)
	if !ok {
		return managed.ExternalDelete{}, errNotBucket
	}

	if hasDeleteAllPolicy(bucket) {
		err := b.deleteAllObjects(ctx, bucket)
		if err != nil {
			return managed.ExternalDelete{}, err
		}
	}

	err := b.deleteS3Bucket(ctx, bucket)
	if err != nil {
		return managed.ExternalDelete{}, err
	}

	b.emitDeletionEvent(bucket)

	return managed.ExternalDelete{}, nil
}
func hasDeleteAllPolicy(bucket *miniov1.Bucket) bool {
	return bucket.Spec.ForProvider.BucketDeletionPolicy == miniov1.DeleteAll
}

func (b *bucketClient) deleteAllObjects(ctx context.Context, bucket *miniov1.Bucket) error {
	log := controllerruntime.LoggerFrom(ctx)
	bucketName := bucket.Status.AtProvider.BucketName

	objectsCh := make(chan minio.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		for object := range b.mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true}) {
			if object.Err != nil {
				log.V(1).Info("warning: cannot list object", "key", object.Key, "error", object.Err)
				continue
			}
			objectsCh <- object
		}
	}()

	bypassGovernance, err := b.isBucketLockEnabled(ctx, bucketName)
	if err != nil {
		log.Error(err, "not able to determine ObjectLock status for bucket", "bucket", bucketName)
	}

	for obj := range b.mc.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{GovernanceBypass: bypassGovernance}) {
		return fmt.Errorf("object %q cannot be removed: %w", obj.ObjectName, obj.Err)
	}
	return nil
}

func (b *bucketClient) isBucketLockEnabled(ctx context.Context, bucketName string) (bool, error) {
	_, mode, _, _, err := b.mc.GetObjectLockConfig(ctx, bucketName)
	if err != nil && err.Error() == "Object Lock configuration does not exist for this bucket" {
		return false, nil
	} else if err != nil {
		return false, err
	}
	// If there's an objectLockConfig it could still be disabled...
	return mode != nil, nil
}

// deleteS3Bucket deletes the bucket.
// NOTE: The removal fails if there are still objects in the bucket.
// This func does not recursively delete all objects beforehand.
func (b *bucketClient) deleteS3Bucket(ctx context.Context, bucket *miniov1.Bucket) error {
	bucketName := bucket.Status.AtProvider.BucketName
	err := b.mc.RemoveBucket(ctx, bucketName)
	return err
}

func (b *bucketClient) emitDeletionEvent(bucket *miniov1.Bucket) {
	b.recorder.Event(bucket, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Deleted",
		Message: "Bucket deleted",
	})
}
