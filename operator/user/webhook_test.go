package user

import (
	"context"
	"encoding/json"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/go-logr/logr"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mockAdminClient struct {
	policies map[string]json.RawMessage
}

func TestValidator_ValidateCreate(t *testing.T) {
	tests := []struct {
		name         string
		obj          *miniov1.User
		wantErr      bool
		wantPolicies map[string]json.RawMessage
	}{
		{
			name: "GivenValidObject_ThenNoError",
			obj: &miniov1.User{
				Spec: miniov1.UserSpec{
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: "test",
						},
					},
				},
			},
		},
		{
			name:    "GivenInvalidObject_ThenError",
			wantErr: true,
			obj:     &miniov1.User{},
		},
		{
			name:    "GivenNotExistingPolicies_ThenError",
			wantErr: true,
			obj: &miniov1.User{
				Spec: miniov1.UserSpec{
					ForProvider: miniov1.UserParameters{
						Policies: []string{
							"foo",
						},
					},
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: "test",
						},
					},
				},
			},
		},
		{
			name: "GivenExistingPolicies_ThenNoError",
			obj: &miniov1.User{
				Spec: miniov1.UserSpec{
					ForProvider: miniov1.UserParameters{
						Policies: []string{
							"foo",
						},
					},
					ResourceSpec: xpv1.ResourceSpec{
						ProviderConfigReference: &xpv1.Reference{
							Name: "test",
						},
					},
				},
			},
			wantPolicies: map[string]json.RawMessage{
				"foo": []byte("foo"),
			},
		},
	}

	for _, tt := range tests {

		getMinioAdminFn = getMockMinioAdmin(tt.wantPolicies)
		getProviderConfigFn = getMockProviderConfig

		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				log: logr.Discard(),
			}
			_, err := v.ValidateCreate(context.TODO(), tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})

		getMinioAdminFn = getMinioAdmin
		getProviderConfigFn = getProviderConfig

	}
}

func TestValidator_ValidateUpdate(t *testing.T) {
	type args struct {
		oldObj *miniov1.User
		newObj *miniov1.User
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantPolicies map[string]json.RawMessage
	}{
		{
			name: "GivenSameObject_ThenNoError",
			args: args{

				oldObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "provider",
							},
						},
					},
				},
				newObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "provider",
							},
						},
					},
				},
			},
		},
		{
			name: "GivenDifferentProviderConfigRef_ThenNoError",
			args: args{
				oldObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "provider",
							},
						},
					},
				},
				newObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "new",
							},
						},
					},
				},
			},
		},
		{
			name:    "GivenDifferentName_ThenError",
			wantErr: true,
			args: args{
				oldObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "provider",
							},
						},
					},
				},
				newObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ForProvider: miniov1.UserParameters{
							UserName: "test",
						},
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "new",
							},
						},
					},
				},
			},
		},
		{
			name:    "GivenNotExistingPolicies_ThenError",
			wantErr: true,
			args: args{
				oldObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "provider",
							},
						},
					},
				},
				newObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ForProvider: miniov1.UserParameters{
							Policies: []string{
								"foo",
							},
						},
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "test",
							},
						},
					},
				},
			},
		},
		{
			name: "GivenExistingPolicies_ThenNoError",
			args: args{
				oldObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "provider",
							},
						},
					},
				},
				newObj: &miniov1.User{
					Spec: miniov1.UserSpec{
						ForProvider: miniov1.UserParameters{
							Policies: []string{
								"foo",
							},
						},
						ResourceSpec: xpv1.ResourceSpec{
							ProviderConfigReference: &xpv1.Reference{
								Name: "test",
							},
						},
					},
				},
			},
			wantPolicies: map[string]json.RawMessage{
				"foo": []byte("foo"),
			},
		},
	}
	for _, tt := range tests {

		getMinioAdminFn = getMockMinioAdmin(tt.wantPolicies)
		getProviderConfigFn = getMockProviderConfig

		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				log: logr.Discard(),
			}
			_, err := v.ValidateUpdate(context.TODO(), tt.args.oldObj, tt.args.newObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})

		getMinioAdminFn = getMinioAdmin
		getProviderConfigFn = getProviderConfig
	}
}

func getMockProviderConfig(context.Context, *miniov1.User, client.Client) (*providerv1.ProviderConfig, error) {
	return &providerv1.ProviderConfig{}, nil
}

func getMockMinioAdmin(policies map[string]json.RawMessage) func(context.Context, client.Client, *providerv1.ProviderConfig) (cannedPolicyLister, error) {
	return func(context.Context, client.Client, *providerv1.ProviderConfig) (cannedPolicyLister, error) {
		return &mockAdminClient{
			policies: policies,
		}, nil
	}
}

func (m *mockAdminClient) ListCannedPolicies(ctx context.Context) (map[string]json.RawMessage, error) {
	return m.policies, nil
}
