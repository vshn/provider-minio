## These are some common variables for Make

PROJECT_ROOT_DIR = .
# TODO: Adjust project meta
PROJECT_NAME ?= go-bootstrap
PROJECT_OWNER ?= vshn

## Variables relevant for building Go
BIN_FILENAME ?= $(PROJECT_NAME)

## Variables relevant for building with Docker
DOCKER_CMD ?= docker
IMG_TAG ?= latest
# Image URL to use all building/pushing image targets
CONTAINER_IMG ?= local.dev/$(PROJECT_OWNER)/$(PROJECT_NAME):$(IMG_TAG)
