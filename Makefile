# Set Shell to bash, otherwise some targets fail with dash/zsh etc.
SHELL := /bin/bash

# Disable built-in rules
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:
.SECONDARY:
.DEFAULT_GOAL := help

# General variables
include Makefile.vars.mk

# Other makefiles
include package/package.mk
include test/local.mk

# Following includes do not print warnings or error if files aren't found
# Optional Documentation module.
-include docs/antora-preview.mk docs/antora-build.mk
# Optional kind module
-include kind/kind.mk

golangci_bin = $(go_bin)/golangci-lint

.PHONY: help
help: ## Show this help
	@grep -E -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: build-bin build-docker ## All-in-one build

.PHONY: build-bin
build-bin: export CGO_ENABLED = 0
build-bin: fmt vet ## Build binary
	@go build -o $(BIN_FILENAME) .

.PHONY: build-docker
build-docker: build-bin ## Build docker image
	$(DOCKER_CMD) build -t $(CONTAINER_IMG) .

.PHONY: test
test: test-go ## All-in-one test

.PHONY: test-go
test-go: ## Run unit tests against code
	go test -race -covermode atomic ./...

.PHONY: fmt
fmt: ## Run 'go fmt' against code
	go fmt ./...

.PHONY: vet
vet: ## Run 'go vet' against code
	go vet ./...

.PHONY: lint
lint: generate fmt golangci-lint git-diff ## All-in-one linting

.PHONY: golangci-lint
golangci-lint: $(golangci_bin) ## Run golangci linters
	$(golangci_bin) run --timeout 5m --out-format colored-line-number ./...

.PHONY: git-diff
git-diff:
	@echo 'Check for uncommitted changes ...'
	git diff --exit-code

.PHONY: generate
generate: ## Generate additional code and artifacts
	@go generate ./...

.PHONY: clean
clean: kind-clean ## Cleans local build artifacts
	docker rmi $(CONTAINER_IMG) || true
	rm -rf docs/node_modules $(docs_out_dir) dist .cache $(WORK_DIR)

$(golangci_bin): | $(go_bin)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go_bin)"

.PHONY: install-crd
install-crd: export KUBECONFIG = $(KIND_KUBECONFIG)
install-crd: generate ## Install CRDs into cluster
	kubectl apply -f package/crds

.PHONY: install-samples
install-samples: export KUBECONFIG = $(KIND_KUBECONFIG)
install-samples: ## Install samples into cluster
	yq ./samples/exoscale*.yaml | kubectl apply -f -

.PHONY: delete-samples
delete-samples: export KUBECONFIG = $(KIND_KUBECONFIG)
delete-samples:
	-yq ./samples/*.yaml | kubectl delete --ignore-not-found --wait=false -f -

.PHONY: run-operator
run-operator: ## Run in Operator mode against your current kube context
	go run . --log-level 1 operator

# Generate webhook certificates.
# This is only relevant when debugging.
# Component-appcat installs a proper certificate for this.
.PHONY: webhook-cert
webhook_key = .work/webhook/tls.key
webhook_cert = .work/webhook/tls.crt
webhook-cert: $(webhook_cert) ## Generate webhook certificates for out-of-cluster debugging in an IDE

$(webhook_key):
	mkdir .work/webhook
	ipsan="" && \
	if [[ $(webhook_service_name) =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$$ ]]; then \
		ipsan=", IP:$(webhook_service_name)"; \
	fi; \
	openssl req -x509 -newkey rsa:4096 -nodes -keyout $@ --noout -days 3650 -subj "/CN=$(webhook_service_name)" -addext "subjectAltName = DNS:$(webhook_service_name)$$ipsan"

$(webhook_cert): $(webhook_key)
	ipsan="" && \
	if [[ $(webhook_service_name) =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$$ ]]; then \
		ipsan=", IP:$(webhook_service_name)"; \
	fi; \
	openssl req -x509 -key $(webhook_key) -nodes -out $@ -days 3650 -subj "/CN=$(webhook_service_name)" -addext "subjectAltName = DNS:$(webhook_service_name)$$ipsan"


.PHONY: webhook-debug
webhook_service_name = host.docker.internal

webhook-debug: $(webhook_cert) ## Creates certificates, patches the webhook registrations and applies everything to the given kube cluster
webhook-debug:
	#
	# kubectl -n syn-appcat scale deployment appcat-controller --replicas 0
	# cabundle=$$(cat .work/webhook/tls.crt | base64) && \
	# HOSTIP=$(webhook_service_name) && \
	# kubectl annotate validatingwebhookconfigurations.admissionregistration.k8s.io appcat-pg-validation cert-manager.io/inject-ca-from- && \
	# kubectl get validatingwebhookconfigurations.admissionregistration.k8s.io appcat-pg-validation -oyaml | \
	# yq e "del(.webhooks[0].clientConfig.service) | .webhooks[0].clientConfig.caBundle |= \"$$cabundle\" | .webhooks[0].clientConfig.url |= \"https://$$HOSTIP:9443/validate-vshn-appcat-vshn-io-v1-vshnpostgresql\"" - | \
	# kubectl apply -f - && \
	# kubectl annotate validatingwebhookconfigurations.admissionregistration.k8s.io appcat-redis-validation cert-manager.io/inject-ca-from- && \
	# kubectl annotate validatingwebhookconfigurations.admissionregistration.k8s.io appcat-pg-validation kubectl.kubernetes.io/last-applied-configuration- && \

