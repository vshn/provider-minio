package minioutil

import (
	"context"
	"net/url"

	"github.com/minio/madmin-go/v3"
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewMinioAdmin returns a new minio admin client that can manage users and IAM.
// It can be used to assign a policy to a usser.
func NewMinioAdmin(ctx context.Context, c client.Client, config *providerv1.ProviderConfig) (*madmin.AdminClient, error) {

	secret, tls, parsed, err := ExtractDataFromProviderConfig(ctx, c, config)
	if err != nil {
		return nil, err
	}

	return madmin.New(parsed.Host, string(secret.Data[MinioIDKey]), string(secret.Data[MinioSecretKey]), tls)
}

// this is the helper function that is used in the NewMinioClient function
func ExtractDataFromProviderConfig(ctx context.Context, c client.Client, config *providerv1.ProviderConfig) (*corev1.Secret, bool, *url.URL, error) {
	secret := &corev1.Secret{}
	key := client.ObjectKey{Name: config.Spec.Credentials.APISecretRef.Name, Namespace: config.Spec.Credentials.APISecretRef.Namespace}
	err := c.Get(ctx, key, secret)
	if err != nil {
		return nil, false, nil, err
	}

	parsed, err := url.Parse(config.Spec.MinioURL)
	if err != nil {
		return nil, false, nil, err
	}

	tls := isTLSEnabled(parsed)

	return secret, tls, parsed, nil
}
