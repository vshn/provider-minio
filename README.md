# provider-minio

[![Build](https://img.shields.io/github/workflow/status/vshn/provider-minio/Test)][build]
![Go version](https://img.shields.io/github/go-mod/go-version/vshn/provider-minio)
[![Version](https://img.shields.io/github/v/release/vshn/provider-minio)][releases]
[![GitHub downloads](https://img.shields.io/github/downloads/vshn/provider-minio/total)][releases]

[build]: https://github.com/vshn/provider-minio/actions?query=workflow%3ATest
[releases]: https://github.com/vshn/provider-minio/releases

Crossplane provider for managing resources on min.io.

Documentation: https://vshn.github.io/provider-minio/provider-minio/

## Local Development

### Requirements

* `docker`
* `go`
* `helm`
* `kubectl`
* `yq`
* `sed` (or `gsed` for Mac)

Some other requirements (e.g. `kind`) will be compiled on-the-fly and put in the local cache dir `.kind` as needed.

### Common make targets

* `make build` to build the binary and docker image
* `make generate` to (re)generate additional code artifacts
* `make test` run test suite
* `make local-install` to install the operator in local cluster
* `make install-samples` to run the provider in local cluster and apply sample manifests
* `make run-operator` to run the code in operator mode against your current kubecontext

See all targets with `make help`

### QuickStart Demonstration

1. Make sure you have a kind cluster running and the config exported
2. `make local-install`

### Kubernetes Webhook Troubleshooting
The provider comes with mutating and validation admission webhook server.

To test and troubleshoot the webhooks on the cluster, simply apply your changes with `kubectl`.

1.  Make sure you have all CRDs and validation webhook registrations installed.
    ```bash
    make install-crd
    kubectl apply -f package/webhook
    ```
2.  To debug the webhook in an IDE, we need to generate certificates:
    ```bash
    make webhook-debug
    # if necessary with another endpoint name, depending on your docker setup
    # if you change the webhook_service_name variable, you need to clean out the old certificates
    make webhook-debug -e webhook_service_name=$HOSTIP
    ```
3.  Start the operator in your IDE with `WEBHOOK_TLS_CERT_DIR` environment set to `.work/webhooks`.

4.  Apply the samples to test the webhooks:
    ```bash
    make install-samples
    ```

### Run operator in debugger

* `make crossplane-setup minio-setup install-crds` to install crossplane and minio in the kind cluster
* `kubectl apply -f samples/_secret.yaml samples/minio.crossplane.io_providerconfig.yaml`
* `EXPORT KUBECONFIG=.work/kind/kind-kubeconfig`
* `go run . --log-level 1 operator`

### Crossplane Provider Mechanics

For detailed information on how Crossplane Provider works from a development perspective check [provider mechanics documentation page](https://kb.vshn.ch/app-catalog/explanations/crossplane_provider_mechanics.html).

### e2e testing with kuttl

Some scenarios are tested with the Kubernetes E2E testing tool [Kuttl](https://kuttl.dev/docs).
Kuttl is basically comparing the installed manifests (usually files named `##-install*.yaml`) with observed objects and compares the desired output (files named `##-assert*.yaml`).

To execute tests, run `make test-e2e` from the root dir.

If a test fails, kuttl leaves the resources in the kind-cluster intact, so you can inspect the resources and events if necessary.
Please note that Kubernetes Events from cluster-scoped resources appear in the `default` namespace only, but `kubectl describe ...` should show you the events.

### Cleaning up e2e tests

`make clean`
