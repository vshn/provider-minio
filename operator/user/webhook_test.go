package user

import (
	"context"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/go-logr/logr"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
)

func TestValidator_ValidateCreate(t *testing.T) {
	tests := []struct {
		name    string
		obj     *miniov1.User
		wantErr bool
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
	}
	for _, tt := range tests {
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
	}
}

func TestValidator_ValidateUpdate(t *testing.T) {
	type args struct {
		oldObj *miniov1.User
		newObj *miniov1.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
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
	}
	for _, tt := range tests {
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
	}
}
