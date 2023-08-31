package user

import (
	"testing"

	"github.com/minio/madmin-go/v3"
	miniov1 "github.com/vshn/provider-minio/apis/minio/v1"
)

func Test_userClient_equalPolicies(t *testing.T) {
	type args struct {
		minioUser madmin.UserInfo
		user      *miniov1.User
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "GivenBothEmpty_ThenTrue",
			want: true,
			args: args{
				minioUser: madmin.UserInfo{},
				user:      &miniov1.User{},
			},
		},
		{
			name: "GivenPolicyOnUserInfo_ThenFalse",
			want: false,
			args: args{
				minioUser: madmin.UserInfo{
					PolicyName: "mypolicy",
				},
				user: &miniov1.User{},
			},
		},
		{
			name: "GivenBothPolicy_ThenTrue",
			want: true,
			args: args{
				minioUser: madmin.UserInfo{
					PolicyName: "mypolicy",
				},
				user: &miniov1.User{
					Spec: miniov1.UserSpec{
						ForProvider: miniov1.UserParameters{
							Policies: []string{
								"mypolicy",
							},
						},
					},
				},
			},
		},
		{
			name: "GivenMoreMinioPolicies_ThenFalse",
			want: false,
			args: args{
				minioUser: madmin.UserInfo{
					PolicyName: "mypolicy,another",
				},
				user: &miniov1.User{
					Spec: miniov1.UserSpec{
						ForProvider: miniov1.UserParameters{
							Policies: []string{
								"mypolicy",
							},
						},
					},
				},
			},
		},
		{
			name: "GivenNoMinioPolicy_ThenFalse",
			want: false,
			args: args{
				minioUser: madmin.UserInfo{},
				user: &miniov1.User{
					Spec: miniov1.UserSpec{
						ForProvider: miniov1.UserParameters{
							Policies: []string{
								"mypolicy",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userClient{}
			if got := u.equalPolicies(tt.args.minioUser, tt.args.user); got != tt.want {
				t.Errorf("userClient.equalPolicies() = %v, want %v", got, tt.want)
			}
		})
	}
}
