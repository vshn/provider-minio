package v1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// +kubebuilder:object:root=true

// A ProviderConfigUsage indicates that a resource is using a ProviderConfig.
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Config-Name",type="string",JSONPath=".providerConfigRef.name"
// +kubebuilder:printcolumn:name="Resource-Kind",type="string",JSONPath=".resourceRef.kind"
// +kubebuilder:printcolumn:name="Resource-Name",type="string",JSONPath=".resourceRef.name"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,provider,minio}
type ProviderConfigUsage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	xpv1.ProviderConfigUsage `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfigUsageList contains a list of ProviderConfigUsage
type ProviderConfigUsageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfigUsage `json:"items"`
}

// ProviderConfigUsage type metadata.
var (
	ProviderConfigUsageKind             = reflect.TypeOf(ProviderConfigUsage{}).Name()
	ProviderConfigUsageGroupKind        = schema.GroupKind{Group: Group, Kind: ProviderConfigUsageKind}.String()
	ProviderConfigUsageKindAPIVersion   = ProviderConfigUsageKind + "." + SchemeGroupVersion.String()
	ProviderConfigUsageGroupVersionKind = SchemeGroupVersion.WithKind(ProviderConfigUsageKind)

	ProviderConfigUsageListKind             = reflect.TypeOf(ProviderConfigUsageList{}).Name()
	ProviderConfigUsageListGroupKind        = schema.GroupKind{Group: Group, Kind: ProviderConfigUsageListKind}.String()
	ProviderConfigUsageListKindAPIVersion   = ProviderConfigUsageListKind + "." + SchemeGroupVersion.String()
	ProviderConfigUsageListGroupVersionKind = SchemeGroupVersion.WithKind(ProviderConfigUsageListKind)
)

func init() {
	SchemeBuilder.Register(&ProviderConfigUsage{}, &ProviderConfigUsageList{})
}
