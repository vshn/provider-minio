## These are some common variables for Make

PROJECT_ROOT_DIR = .
PROJECT_NAME ?= provider-minio
PROJECT_OWNER ?= vshn

WORK_DIR = $(PWD)/.work

## BUILD:go
BIN_FILENAME ?= $(PROJECT_NAME)
go_bin ?= $(WORK_DIR)/bin
$(go_bin):
	@mkdir -p $@

## BUILD:docker
DOCKER_CMD ?= docker
CONTAINER_REGISTRY ?= ghcr.io
UPBOUND_CONTAINER_REGISTRY ?= xpkg.upbound.io

GIT_TAG = $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
IMG_TAG = $(subst /,_,$(GIT_TAG))
# Image URL to use all building/pushing image targets
CONTAINER_IMG ?= $(CONTAINER_REGISTRY)/$(PROJECT_OWNER)/$(PROJECT_NAME)/controller:$(IMG_TAG)
LOCAL_PACKAGE_IMG = localhost:15000/$(PROJECT_OWNER)/$(PROJECT_NAME)/package:$(IMG_TAG)
GHCR_PACKAGE_IMG ?= $(CONTAINER_REGISTRY)/$(PROJECT_OWNER)/$(PROJECT_NAME)/provider:$(IMG_TAG)
UPBOUND_PACKAGE_IMG ?= $(UPBOUND_CONTAINER_REGISTRY)/$(PROJECT_OWNER)/$(PROJECT_NAME):$(IMG_TAG)

## KIND:setup

# https://hub.docker.com/r/kindest/node/tags
KIND_NODE_VERSION ?= v1.26.6
KIND_IMAGE ?= docker.io/kindest/node:$(KIND_NODE_VERSION)
KIND ?= go run sigs.k8s.io/kind
KIND_KUBECONFIG ?= $(kind_dir)/kind-kubeconfig
KIND_CLUSTER ?= $(PROJECT_NAME)
