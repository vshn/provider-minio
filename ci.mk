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
DOCKER_IMAGE_GOOS = linux
DOCKER_IMAGE_GOARCH = amd64

.PHONY: docker-build
docker-build:
	env CGO_ENABLED=0 GOOS=$(DOCKER_IMAGE_GOOS) GOARCH=$(DOCKER_IMAGE_GOARCH) \
		go build -o ${BIN_FILENAME}
	docker build --platform $(DOCKER_IMAGE_GOOS)/$(DOCKER_IMAGE_GOARCH) -t ${IMG} .

.PHONY: docker-build-branchtag
docker-build-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
docker-build-branchtag: docker-build ## Build docker image with current branch name

.PHONY: docker-push
docker-push: docker-build ## Push docker image with the manager.
	docker push ${IMG}

.PHONY: docker-push-branchtag
docker-push-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
docker-push-branchtag: docker-build-branchtag docker-push ## Push docker image with current branch name

.PHONY: package-build
package-build: docker-build
	rm -f package/*.xpkg
	go run github.com/crossplane/crossplane/cmd/crank@v1.16.0 xpkg build -f package --verbose --embed-runtime-image=${IMG} -o package/package.xpkg

.PHONY: package-push
package-push: package-build
	go run github.com/crossplane/crossplane/cmd/crank@v1.16.0 xpkg push -f package/package.xpkg ${IMG} --verbose

.PHONY: package-build-branchtag
package-build-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
package-build-branchtag: docker-build-branchtag package-build

.PHONY: package-push-package-branchtag
package-push-package-branchtag: export IMG_TAG=$(shell git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')
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
