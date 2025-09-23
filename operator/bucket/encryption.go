package bucket

import (
	"fmt"

	"github.com/minio/minio-go/v7/pkg/sse"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
)

// Set up an SSE encryption struct based on this buckets encryption specification.
// Returns true and an SSE config struct if encryption is set up.
// Returns false if encryption is not enabled and/or an error if it is improperly set up.
func createSseConfig(bucket *miniov1.Bucket) (bool, sse.Configuration, error) {
	sseConfig := sse.Configuration{}

	encSpec := bucket.Spec.ForProvider.Encryption
	switch encSpec.Type {
	case "sse-s3":
		sseConfig = *sse.NewConfigurationSSES3()
	case "sse-kms":
		if encSpec.KmsId == "" {
			return false, sse.Configuration{}, fmt.Errorf("must set KMS key ID if encryption type is set to 'sse-kms'")
		}
		sseConfig = *sse.NewConfigurationSSEKMS(encSpec.KmsId)
	default:
		return false, sse.Configuration{}, nil
	}

	return true, sseConfig, nil
}
