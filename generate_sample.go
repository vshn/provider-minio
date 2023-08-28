//go:build generate

// Clean samples dir
//go:generate rm -rf ./samples/*

// Generate sample files
//go:generate go run generate_sample.go ./samples

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vshn/provider-minio/apis"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	"github.com/vshn/provider-minio/operator/minioutil"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	serializerjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

var (
	scheme = runtime.NewScheme()
)

func main() {
	failIfError(apis.AddToScheme(scheme))
	generateBucketSample()
	generateProviderConfigSample()
	generateSecretSample()
}

func generateBucketSample() {
	spec := newBucketSample()
	serialize(spec, true)
}

func newBucketSample() *miniov1.Bucket {
	return &miniov1.Bucket{
		TypeMeta: metav1.TypeMeta{
			APIVersion: miniov1.BucketGroupVersionKind.GroupVersion().String(),
			Kind:       miniov1.BucketKind,
		},
		ObjectMeta: metav1.ObjectMeta{Name: "bucket-local-dev"},
		Spec: miniov1.BucketSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{Name: "provider-config"},
			},
			ForProvider: miniov1.BucketParameters{
				Region: "us-east-1",
			},
		},
	}
}

func generateProviderConfigSample() {
	spec := newProviderConfigSample()
	serialize(spec, true)
}

func newProviderConfigSample() *providerv1.ProviderConfig {
	return &providerv1.ProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: providerv1.ProviderConfigGroupVersionKind.GroupVersion().String(),
			Kind:       providerv1.ProviderConfigKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "provider-config"},
		Spec: providerv1.ProviderConfigSpec{
			MinioURL: "http://minio.127.0.0.1.nip.io:8088/",
			Credentials: providerv1.ProviderCredentials{
				Source: xpv1.CredentialsSourceInjectedIdentity,
				APISecretRef: corev1.SecretReference{
					Name:      "minio-secret",
					Namespace: "crossplane-system",
				},
			},
		},
	}
}

func generateSecretSample() {
	spec := newSecretSample()
	serialize(spec, true)
}

func newSecretSample() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minio-secret",
			Namespace: "crossplane-system",
		},
		Data: map[string][]byte{
			minioutil.MinioIDKey:     []byte("minioadmin"),
			minioutil.MinioSecretKey: []byte("minioadmin"),
		},
	}
}

func failIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func serialize(object runtime.Object, useYaml bool) {
	gvk := object.GetObjectKind().GroupVersionKind()
	fileExt := "json"
	if useYaml {
		fileExt = "yaml"
	}
	fileName := fmt.Sprintf("%s_%s.%s", strings.ToLower(gvk.Group), strings.ToLower(gvk.Kind), fileExt)
	f := prepareFile(fileName)
	err := serializerjson.NewSerializerWithOptions(serializerjson.DefaultMetaFactory, scheme, scheme, serializerjson.SerializerOptions{Yaml: useYaml, Pretty: true}).Encode(object, f)
	failIfError(err)
}

func prepareFile(file string) io.Writer {
	dir := os.Args[1]
	err := os.MkdirAll(os.Args[1], 0775)
	failIfError(err)
	f, err := os.Create(filepath.Join(dir, file))
	failIfError(err)
	return f
}
