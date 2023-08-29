package bucket

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/minio/minio-go/v7"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func (b *bucketClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	log := controllerruntime.LoggerFrom(ctx)
	log.V(1).Info("creating resource")

	bucket, ok := mg.(*miniov1.Bucket)
	if !ok {
		return managed.ExternalCreation{}, errNotBucket
	}

	err := b.createS3Bucket(ctx, bucket)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	b.setLock(bucket)

	return managed.ExternalCreation{}, b.emitCreationEvent(bucket)
}

// createS3Bucket creates a new bucket and sets the name in the status.
// If the bucket already exists, and we have permissions to access it, no error is returned and the name is set in the status.
// If the bucket exists, but we don't own it, an error is returned.
func (b *bucketClient) createS3Bucket(ctx context.Context, bucket *miniov1.Bucket) error {
	bucketName := bucket.GetBucketName()
	err := b.mc.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: bucket.Spec.ForProvider.Region})

	if err != nil {
		// Check to see if we already own this bucket (which happens if we run this twice)
		exists, errBucketExists := b.mc.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil
		}
		// someone else might have created the bucket
		return err

	}
	return nil
}

// setLock sets an annotation that tells the Observe func that we have successfully created the bucket.
// Without it, another resource that has the same bucket name might "adopt" the same bucket, causing 2 resources managing 1 bucket.
func (b *bucketClient) setLock(bucket *miniov1.Bucket) {
	if bucket.Annotations == nil {
		bucket.Annotations = map[string]string{}
	}
	bucket.Annotations[lockAnnotation] = "claimed"

}

func (b *bucketClient) emitCreationEvent(bucket *miniov1.Bucket) error {
	b.recorder.Event(bucket, event.Event{
		Type:    event.TypeNormal,
		Reason:  "Created",
		Message: "Bucket successfully created",
	})
	return nil
}
