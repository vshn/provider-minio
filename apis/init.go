// Package apis contains Kubernetes API for the Template provider.
package apis

import (
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	"k8s.io/apimachinery/pkg/runtime"

	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
)

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		miniov1.SchemeBuilder.AddToScheme,
		providerv1.SchemeBuilder.AddToScheme,
	)
}

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
