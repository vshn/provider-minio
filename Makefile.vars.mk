## These are some common variables for Make

PROJECT_ROOT_DIR = .
# TODO: Adjust project meta
PROJECT_NAME ?= go-bootstrap
PROJECT_OWNER ?= vshn

## BUILD:go
BIN_FILENAME ?= $(PROJECT_NAME)

## BUILD:docker
DOCKER_CMD ?= docker

IMG_TAG ?= latest
# Image URL to use all building/pushing image targets
CONTAINER_IMG ?= local.dev/$(PROJECT_OWNER)/$(PROJECT_NAME):$(IMG_TAG)


## KIND:setup

# https://hub.docker.com/r/kindest/node/tags
KIND_NODE_VERSION ?= v1.24.0
KIND_IMAGE ?= docker.io/kindest/node:$(KIND_NODE_VERSION)
KIND ?= go run sigs.k8s.io/kind
KIND_KUBECONFIG ?= $(kind_dir)/kind-kubeconfig-$(KIND_NODE_VERSION)
KIND_CLUSTER ?= $(PROJECT_NAME)-$(KIND_NODE_VERSION)
