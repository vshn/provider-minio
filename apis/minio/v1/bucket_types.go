package v1

import (
	"reflect"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	SchemeBuilder.Register(&Bucket{}, &BucketList{})
}

const (
	// DeleteIfEmpty only deletes the bucket if the bucket is empty.
	DeleteIfEmpty BucketDeletionPolicy = "DeleteIfEmpty"
	// DeleteAll recursively deletes all objects in the bucket and then removes it.
	DeleteAll BucketDeletionPolicy = "DeleteAll"
)

// BucketDeletionPolicy determines how buckets should be deleted when a Bucket is deleted.
type BucketDeletionPolicy string

// We can't have this here, because ironically the generator breaks if this throws and error...
// var _ resource.Managed = &Bucket{}
// var _ runtime.Object = &Bucket{}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.endpointURL"
// +kubebuilder:printcolumn:name="Bucket Name",type="string",JSONPath=".status.atProvider.bucketName"
// +kubebuilder:printcolumn:name="Region",type="string",JSONPath=".spec.forProvider.region"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,minio}
// +kubebuilder:webhook:verbs=create;update,path=/validate-minio-crossplane-io-v1-bucket,mutating=false,failurePolicy=fail,groups=minio.crossplane.io,resources=buckets,versions=v1,name=buckets.minio.crossplane.io,sideEffects=None,admissionReviewVersions=v1

type Bucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BucketSpec   `json:"spec"`
	Status BucketStatus `json:"status,omitempty"`
}

type BucketSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       BucketParameters `json:"forProvider,omitempty"`
}

type BucketStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	Endpoint            string               `json:"endpoint,omitempty"`
	EndpointURL         string               `json:"endpointURL,omitempty"`
	AtProvider          BucketProviderStatus `json:"atProvider,omitempty"`
}

type BucketParameters struct {
	// BucketName is the name of the bucket to create.
	// Defaults to `metadata.name` if unset.
	// Cannot be changed after bucket is created.
	// Name must be acceptable by the S3 protocol, which follows RFC 1123.
	// Be aware that S3 providers may require a unique name across the platform or zone.
	BucketName string `json:"bucketName,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:default="us-east-1"

	// Region is the name of the region where the bucket shall be created.
	// The region must be available in the S3 endpoint.
	// Cannot be changed after bucket is created.
	Region string `json:"region,omitempty"`

	// BucketDeletionPolicy determines how buckets should be deleted when Bucket is deleted.
	//  `DeleteIfEmpty` only deletes the bucket if the bucket is empty.
	//  `DeleteAll` recursively deletes all objects in the bucket and then removes it.
	// To skip deletion of the bucket (orphan it) set `spec.deletionPolicy=Orphan`.
	BucketDeletionPolicy BucketDeletionPolicy `json:"bucketDeletionPolicy,omitempty"`

	// Policy is a raw S3 bucket policy.
	// Please consult https://min.io/docs/minio/linux/administration/identity-access-management/policy-based-access-control.html for more details about the policy.
	Policy *string `json:"policy,omitempty"`
}

type BucketProviderStatus struct {
	// BucketName is the name of the actual bucket.
	BucketName string `json:"bucketName,omitempty"`
}

// +kubebuilder:object:root=true

type BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bucket `json:"items"`
}

// Dummy type metadata.
var (
	BucketKind             = reflect.TypeOf(Bucket{}).Name()
	BucketGroupKind        = schema.GroupKind{Group: Group, Kind: BucketKind}.String()
	BucketKindAPIVersion   = BucketKind + "." + SchemeGroupVersion.String()
	BucketGroupVersionKind = SchemeGroupVersion.WithKind(BucketKind)
)

// GetBucketName returns the spec.forProvider.bucketName if given, otherwise defaults to metadata.name.
func (in *Bucket) GetBucketName() string {
	if in.Spec.ForProvider.BucketName == "" {
		return in.Name
	}
	return in.Spec.ForProvider.BucketName
}
