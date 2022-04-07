
antora_preview_version ?= 3.0.1.1
antora_preview_cmd ?= $(DOCKER_CMD) run --rm --publish 2020:2020 --volume "${PWD}":/preview/antora docker.io/vshn/antora-preview:$(antora_preview_version) --style=vshn

.PHONY: docs-preview
docs-preview: ## Preview the documentation
	$(antora_preview_cmd)
