kind_dir ?= $(PWD)/.kind
kind_bin = $(go_bin)/kind

# Prepare kind binary
$(kind_bin): export GOOS = $(shell go env GOOS)
$(kind_bin): export GOARCH = $(shell go env GOARCH)
$(kind_bin): export GOBIN = $(go_bin)
$(kind_bin): | $(go_bin)
	go install sigs.k8s.io/kind@latest

mirror_sentinel = $(kind_dir)/mirror_sentinel

.PHONY: kind
kind: export KUBECONFIG = $(KIND_KUBECONFIG)
kind: kind-setup-ingress kind-load-image ## All-in-one kind target

.PHONY: kind-setup
kind-setup: export KUBECONFIG = $(KIND_KUBECONFIG)
kind-setup: $(KIND_KUBECONFIG) ## Creates the kind cluster

.PHONY: kind-setup-ingress
kind-setup-ingress: export KUBECONFIG = $(KIND_KUBECONFIG)
kind-setup-ingress: kind-setup ## Install NGINX as ingress controller onto kind cluster (localhost:8088)
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

.PHONY: kind-load-image
# We fix the arch to linux/amd64 since kind runs in amd64 even on Mac/arm.
kind-load-image: export GOOS = linux
kind-load-image: export GOARCH = amd64
kind-load-image: kind-setup docker-build ## Load the container image onto kind cluster
	@$(kind_bin) load docker-image --name $(KIND_CLUSTER) $(IMG)

.PHONY: kind-clean
kind-clean: export KUBECONFIG = $(KIND_KUBECONFIG)
kind-clean: ## Removes the kind Cluster
	@$(kind_bin) delete cluster --name $(KIND_CLUSTER) || true
	docker rm -f kind-registry
	rm -rf $(kind_dir) $(kind_bin)

$(KIND_KUBECONFIG): export KUBECONFIG = $(KIND_KUBECONFIG)
$(KIND_KUBECONFIG): $(kind_bin)
	$(kind_bin) create cluster \
		--name $(KIND_CLUSTER) \
		--image $(KIND_IMAGE) \
		--config kind/config.yaml
	@kubectl version
	@kubectl cluster-info
	@kubectl config use-context kind-$(KIND_CLUSTER)
	@echo =======
	@echo "Setup finished. To interact with the local dev cluster, set the KUBECONFIG environment variable as follows:"
	@echo "export KUBECONFIG=$$(realpath "$(KIND_KUBECONFIG)")"
	@echo =======

.PHONY: mirror-setup
mirror-setup: $(mirror_sentinel) ## Installs an image registry required for the package image in kind cluster.

$(mirror_sentinel): export KUBECONFIG = $(KIND_KUBECONFIG)
$(mirror_sentinel):

	REGISTRY_DIR="/etc/containerd/certs.d/registry.registry-system.svc.cluster.local:5000" && \
	REGISTRY_HOST='[host."http://localhost:30500"]' && \
	for node in $$(kind get nodes -n $(KIND_CLUSTER)); do \
		echo $$node ; \
	  docker exec "$${node}" mkdir -p "$${REGISTRY_DIR}" ; \
	  echo "$${REGISTRY_HOST}" | docker exec -i "$${node}" cp /dev/stdin "$${REGISTRY_DIR}/hosts.toml" ; \
	done

	@touch $@
