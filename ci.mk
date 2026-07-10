# Image URL to use all building/pushing image targets
IMG_TAG ?= latest
APP_NAME ?= provider-minio
ORG ?= vshn
IMG_REPO ?= ghcr.io
IMG ?= $(IMG_REPO)/$(ORG)/$(APP_NAME):$(IMG_TAG)
DOCKER_CMD ?= docker

# Upbound push config
UPBOUND_CONTAINER_REGISTRY ?= xpkg.upbound.io
UPBOUND_PACKAGE_IMG ?= $(UPBOUND_CONTAINER_REGISTRY)/$(ORG)/$(APP_NAME):$(IMG_TAG)

# For alpine image it is required the following env before building the application
# DOCKER_IMAGE_GOARCH is a space-separated list; one runtime image is built per
# architecture and the package is pushed as a multi-arch index.
DOCKER_IMAGE_GOOS = linux
DOCKER_IMAGE_GOARCH ?= amd64 arm64

# Comma-join helper for crank's -f flag
empty :=
space := $(empty) $(empty)
comma := ,
XPKG_FILES = $(patsubst %,package/package-%.xpkg,$(DOCKER_IMAGE_GOARCH))

.PHONY: docker-build
docker-build:
	@for arch in $(DOCKER_IMAGE_GOARCH); do \
		env CGO_ENABLED=0 GOOS=$(DOCKER_IMAGE_GOOS) GOARCH=$$arch \
			go build -o ${BIN_FILENAME} && \
		$(DOCKER_CMD) build --platform $(DOCKER_IMAGE_GOOS)/$$arch -t ${IMG}-$$arch . \
		|| exit 1; \
	done
	@# Keep the un-suffixed tag for single-image consumers (kind-load-image);
	@# points at the first architecture in the list (amd64 by default).
	$(DOCKER_CMD) tag ${IMG}-$(firstword $(DOCKER_IMAGE_GOARCH)) ${IMG}

.PHONY: docker-build-branchtag
docker-build-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
docker-build-branchtag: docker-build ## Build docker image with current branch name

.PHONY: docker-push
docker-push: docker-build ## Push docker image with the manager.
	$(DOCKER_CMD) push ${IMG}

.PHONY: docker-push-branchtag
docker-push-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
docker-push-branchtag: docker-build-branchtag docker-push ## Push docker image with current branch name

.PHONY: package-build
package-build: docker-build
	rm -f package/*.xpkg
	@for arch in $(DOCKER_IMAGE_GOARCH); do \
		go run github.com/crossplane/crossplane/cmd/crank@v1.16.0 xpkg build -f package --verbose --embed-runtime-image=${IMG}-$$arch -o package/package-$$arch.xpkg \
		|| exit 1; \
	done

.PHONY: package-push
package-push: package-build
	go run github.com/crossplane/crossplane/cmd/crank@v1.16.0 xpkg push -f $(subst $(space),$(comma),$(XPKG_FILES)) ${IMG} --verbose

.PHONY: package-build-branchtag
package-build-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
package-build-branchtag: docker-build-branchtag package-build

.PHONY: package-push-package-branchtag
package-push-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
package-push-branchtag: package-build-branchtag package-push

.PHONY: docker-build-local
docker-build-local: export IMG_REPO=localhost:5000
docker-build-local:
	$(MAKE) docker-build

.PHONY: package-build-local
package-build-local: export IMG_REPO=localhost:5000
package-build-local: docker-build-local package-build

.PHONY: package-push-local
package-push-local: export IMG_REPO=localhost:5000
package-push-local: package-build-local package-push
