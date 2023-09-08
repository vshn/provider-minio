
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
package_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

crossplane_bin = $(go_bin)/kubectl-crossplane
up_bin = $(go_bin)/up

# Build kubectl-crossplane plugin
$(crossplane_bin):export GOBIN = $(go_bin)
$(crossplane_bin): | $(go_bin)
	go install github.com/crossplane/crossplane/cmd/crank@latest
	@mv $(go_bin)/crank $@

# Install up plugin
$(up_bin):export GOBIN = $(go_bin)
$(up_bin):export VERSION=v0.15.0
$(up_bin): | $(go_bin)
	curl -sL "https://cli.upbound.io" | sh
	@mv up $@
	$(up_bin) --version

.PHONY: package
package: ## All-in-one packaging and releasing
package: package-push

.PHONY: package-provider-local
package-provider-local: export CONTROLLER_IMG = $(CONTAINER_IMG)
package-provider-local: $(crossplane_bin) generate ## Build Crossplane package for local installation in kind-cluster
	@rm -rf package/*.xpkg
	@yq e '.spec.controller.image=strenv(CONTROLLER_IMG)' $(package_dir)/crossplane.yaml.template > $(package_dir)/crossplane.yaml
	@$(crossplane_bin) build provider -f $(package_dir)
	@echo Package file: $$(ls $(package_dir)/*.xpkg)

.PHONY: package-provider
package-provider: export CONTROLLER_IMG = $(CONTAINER_IMG)
package-provider: $(up_bin) generate build-docker ## Build Crossplane package for Upbound Marketplace
	@rm -rf package/*.xpkg
	@yq e 'del(.spec)' $(package_dir)/crossplane.yaml.template > $(package_dir)/crossplane.yaml
	$(up_bin) xpkg build -f $(package_dir) -o $(package_dir)/provider-minio.xpkg --controller=$(CONTROLLER_IMG)

.PHONY: .local-package-push
.local-package-push: pkg_file = $(shell ls $(package_dir)/*.xpkg)
.local-package-push: $(crossplane_bin) package-provider-local
	$(crossplane_bin) push provider -f $(pkg_file) $(LOCAL_PACKAGE_IMG)

.PHONY: .ghcr-package-push
.ghcr-package-push: pkg_file = $(package_dir)/provider-minio.xpkg
.ghcr-package-push: $(crossplane_bin) package-provider
	$(crossplane_bin) push provider -f $(pkg_file) $(GHCR_PACKAGE_IMG)

.PHONY: .upbound-package-push
.upbound-package-push: pkg_file = $(package_dir)/provider-minio.xpkg
.upbound-package-push: package-provider
	$(up_bin) xpkg push -f $(pkg_file) $(UPBOUND_PACKAGE_IMG)

.PHONY: package-push
package-push: .ghcr-package-push ## Push Crossplane package to container registries

.PHONY: .package-clean
.package-clean:
	rm -f $(crossplane_bin) $(up_bin) package/*.xpkg
