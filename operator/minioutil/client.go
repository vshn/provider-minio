package minioutil

import (
	"context"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MinioIDKey     = "AWS_ACCESS_KEY_ID"
	MinioSecretKey = "AWS_SECRET_ACCESS_KEY"
)

// NewMinioClient returns a new minio client according to the given provider config.
func NewMinioClient(ctx context.Context, c client.Client, config *providerv1.ProviderConfig) (*minio.Client, error) {
	secret := &corev1.Secret{}
	key := client.ObjectKey{Name: config.Spec.Credentials.APISecretRef.Name, Namespace: config.Spec.Credentials.APISecretRef.Namespace}
	err := c.Get(ctx, key, secret)
	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(config.Spec.MinioURL)
	if err != nil {
		return nil, err
	}

	return minio.New(parsed.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(string(secret.Data[MinioIDKey]), string(secret.Data[MinioSecretKey]), ""),
		Secure: isTLSEnabled(parsed),
	})

}

// isTLSEnabled returns false if the scheme is explicitly set to `http` or `HTTP`
func isTLSEnabled(u *url.URL) bool {
	return !strings.EqualFold(u.Scheme, "http")
}
