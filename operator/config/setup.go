package config

import (
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/providerconfig"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	ctrl "sigs.k8s.io/controller-runtime"

	providerv1 "github.com/vshn/provider-minio/apis/provider/v1"
)

// SetupController adds a controller that reconciles ProviderConfigs and tracks
// their current usage.
func SetupController(mgr ctrl.Manager) error {
	name := providerconfig.ControllerName(providerv1.ProviderConfigGroupKind)
	recorder := event.NewAPIRecorder(mgr.GetEventRecorderFor(name))

	of := resource.ProviderConfigKinds{
		Config:    providerv1.ProviderConfigGroupVersionKind,
		UsageList: providerv1.ProviderConfigUsageListGroupVersionKind,
	}

	r := providerconfig.NewReconciler(mgr, of,
		providerconfig.WithLogger(logging.NewLogrLogger(mgr.GetLogger().WithValues("controller", name))),
		providerconfig.WithRecorder(recorder))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&providerv1.ProviderConfig{}).
		Watches(&providerv1.ProviderConfigUsage{}, &resource.EnqueueRequestForProviderConfig{}).
		Complete(r)
}
