docs_out_dir := ./.public

docker_opts ?= --rm --tty --user "$$(id -u)"

antora_build_version ?= 2.3.3
antora_preview_version ?= 2.3.14
antora_preview_cmd ?= $(DOCKER_CMD) run --rm --publish 2020:2020 --volume "${PWD}":/antora docker.io/vshn/antora-preview:$(antora_preview_version) --style=vshn --antora=docs

antora_cmd ?= $(DOCKER_CMD) run $(docker_opts) --volume "$${PWD}":/antora docker.io/vshn/antora:$(antora_build_version)
antora_opts ?= --cache-dir=.cache/antora

.PHONY: docs
docs: docs-html ## All-in-one docs build

.PHONY: docs-preview
docs-preview: ## Preview the documentation
	$(antora_preview_cmd)

.PHONY: docs-html
docs-html: $(docs_out_dir)/index.html ## Generate HTML version of documentation with Antora, output at ./.public
	@touch $(docs_out_dir)/.nojekyll

.PHONY: docs-publish
docs-publish: node_modles docs-html ## Publishes the Antora documentation on Github Pages
	npm run deploy

$(docs_out_dir)/index.html:
	$(antora_cmd) $(antora_opts) docs/antora-playbook.yml

node_modules:
	npm install
