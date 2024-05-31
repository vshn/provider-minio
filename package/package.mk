

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
package_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

crossplane_bin = $(go_bin)/kubectl-crossplane

# Build kubectl-crossplane plugin
$(crossplane_bin):export GOBIN = $(go_bin)
$(crossplane_bin): | $(go_bin)
	go install github.com/crossplane/crossplane/cmd/crank@latest
	@mv $(go_bin)/crank $@

.PHONY: package
package: ## All-in-one packaging and releasing
package: package-push

.PHONY: package-provider-local
package-provider-local: export CONTROLLER_IMG = $(CONTAINER_IMG)
package-provider-local: $(crossplane_bin) generate ## Build Crossplane package for local installation in kind-cluster
	@rm -rf package/*.xpkg
	@yq e '.spec.controller.image=strenv(CONTROLLER_IMG)' $(package_dir)/crossplane.yaml.template > $(package_dir)/crossplane.yaml
	@$(crossplane_bin) xpkg build -f $(package_dir)
	@echo Package file: $$(ls $(package_dir)/*.xpkg)

.PHONY: package-provider
package-provider: export CONTROLLER_IMG = $(CONTAINER_IMG)
package-provider: $(crossplane_bin) generate build-docker ## Build Crossplane package for Upbound Marketplace
	@rm -rf package/*.xpkg
	@yq e 'del(.spec)' $(package_dir)/crossplane.yaml.template > $(package_dir)/crossplane.yaml
	$(crossplane_bin) xpkg build --package-root=$(package_dir) --embed-runtime-image=$(CONTROLLER_IMG) -o $(package_dir)/provider-minio.xpkg

.PHONY: .local-package-push
.local-package-push: pkg_file = $(shell ls $(package_dir)/*.xpkg)
.local-package-push: $(crossplane_bin) package-provider-local
	$(crossplane_bin) xpkg push -f $(pkg_file) $(LOCAL_PACKAGE_IMG)

.PHONY: .ghcr-package-push
.ghcr-package-push: pkg_file = $(package_dir)/provider-minio.xpkg
.ghcr-package-push: $(crossplane_bin) package-provider
	$(crossplane_bin) xpkg push -f $(pkg_file) $(GHCR_PACKAGE_IMG)

.PHONY: .upbound-package-push
.upbound-package-push: pkg_file = $(package_dir)/provider-minio.xpkg
.upbound-package-push: package-provider
	$(crossplane_bin) xpkg push -f $(pkg_file) $(UPBOUND_PACKAGE_IMG)

.PHONY: package-push
package-push: .ghcr-package-push .upbound-package-push ## Push Crossplane package to container registries

.PHONY: .package-clean
.package-clean:
	rm -f $(crossplane_bin) package/*.xpkg
