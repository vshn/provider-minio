//go:build generate

// Remove existing manifests
//go:generate rm -rf ../package/crds ../package/webhook

// Generate deepcopy methodsets and CRD manifests
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile=../.github/boilerplate.go.txt paths=./... crd:crdVersions=v1 output:artifacts:config=../package/crds

// Generate webhook manifests
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen webhook paths=./... output:artifacts:config=../package/webhook

// Generate crossplane-runtime methodsets (resource.Claim, etc)
//go:generate go run -tags generate github.com/crossplane/crossplane-tools/cmd/angryjet generate-methodsets --header-file=../.github/boilerplate.go.txt ./...

package apis
