package bucket

import (
	"testing"

	"github.com/minio/minio-go/v7/pkg/sse"
	"github.com/stretchr/testify/assert"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
)

func TestEncryptionSetup(t *testing.T) {
	bucket := miniov1.Bucket{
		Spec: miniov1.BucketSpec{
			ForProvider: miniov1.BucketParameters{
				Encryption: miniov1.BucketEncryption{},
			},
		},
	}

	// Given no encryption, expect no SSE config
	check, config, err := createSseConfig(&bucket)
	assert.False(t, check)
	assert.NoError(t, err)
	assert.Equal(t, sse.Configuration{}, config)

	// Given S3 encryption, expect SSE config
	s3Bucket := bucket.DeepCopy()
	s3Bucket.Spec.ForProvider.Encryption.Type = "sse-s3"
	check, config, err = createSseConfig(s3Bucket)
	assert.True(t, check)
	assert.NoError(t, err)
	assert.Equal(t, "AES256", config.Rules[0].Apply.SSEAlgorithm)

	// Given KMS encryption with KMS ID, expect SSE config
	kmsBucket := bucket.DeepCopy()
	kmsBucket.Spec.ForProvider.Encryption.Type = "sse-kms"
	kmsBucket.Spec.ForProvider.Encryption.KmsId = "most-guarded-secret"
	check, config, err = createSseConfig(kmsBucket)
	assert.True(t, check)
	assert.NoError(t, err)
	assert.Equal(t, "aws:kms", config.Rules[0].Apply.SSEAlgorithm)
	assert.Equal(t, "most-guarded-secret", config.Rules[0].Apply.KmsMasterKeyID)

	// Given KMS encryption but no KMS ID, expect error
	badBucket := bucket.DeepCopy()
	badBucket.Spec.ForProvider.Encryption.Type = "sse-kms"
	check, config, err = createSseConfig(badBucket)
	assert.False(t, check)
	assert.Error(t, err)
	assert.Equal(t, sse.Configuration{}, config)
}
