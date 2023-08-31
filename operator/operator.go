package operator

import (
	"github.com/vshn/provider-minio/operator/bucket"
	"github.com/vshn/provider-minio/operator/policy"
	"github.com/vshn/provider-minio/operator/user"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SetupControllers creates all controllers and adds them to the supplied manager.
func SetupControllers(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		bucket.SetupController,
		user.SetupController,
		policy.SetupController,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhooks creates all webhooks and adds them to the supplied manager.
func SetupWebhooks(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		bucket.SetupWebhook,
		user.SetupWebhook,
		policy.SetupWebhook,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
