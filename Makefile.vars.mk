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

GIT_TAG = $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
IMG_TAG = $(subst /,_,$(GIT_TAG))
# Image URL to use all building/pushing image targets
LOCAL_PACKAGE_IMG = localhost:5000/$(PROJECT_OWNER)/$(PROJECT_NAME)/package:$(IMG_TAG)

## KIND:setup

# https://hub.docker.com/r/kindest/node/tags
KIND_NODE_VERSION ?= v1.28.9
KIND_IMAGE ?= docker.io/kindest/node:$(KIND_NODE_VERSION)
KIND ?= go run sigs.k8s.io/kind
KIND_KUBECONFIG ?= $(kind_dir)/kind-kubeconfig
KIND_CLUSTER ?= $(PROJECT_NAME)
