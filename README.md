# go-bootstrap

[![Build](https://img.shields.io/github/workflow/status/vshn/go-bootstrap/Test)][build]
![Go version](https://img.shields.io/github/go-mod/go-version/vshn/go-bootstrap)
[![Version](https://img.shields.io/github/v/release/vshn/go-bootstrap)][releases]
[![Maintainability](https://img.shields.io/codeclimate/maintainability/vshn/go-bootstrap)][codeclimate]
[![Coverage](https://img.shields.io/codeclimate/coverage/vshn/go-bootstrap)][codeclimate]
[![GitHub downloads](https://img.shields.io/github/downloads/vshn/go-bootstrap/total)][releases]

[build]: https://github.com/vshn/go-bootstrap/actions?query=workflow%3ATest
[releases]: https://github.com/vshn/go-bootstrap/releases
[codeclimate]: https://codeclimate.com/github/vshn/go-bootstrap

Template repository for common Go setups

## Features

* GitHub Workflows
  - Build (Go & Docker image)
  - Test (including CodeClimate)
  - Lint (Go)
  - Release (Goreleaser & Changelog generator)

* GitHub issue templates
  - PR template
  - Issue templates using GitHub issue forms

* Goreleaser
  - Go build for `amd64`, `armv8`
  - Docker build for `latest` and `vx.y.z` tags
  - Push Docker image to GitHub's registry `ghcr.io`

* Antora documentation
  - Build default documentation with VSHN styling
  - Publish to GitHub Pages by default (opt-out)
  - Automated with GitHub workflows to build in `master` branch and (pre-)releases.
  - Available `make` targets are prefixed with `docs-`

* Local Kubernetes environment
  - Setup Kubernetes-In-Docker (kind)
  - Prepares a kubeconfig file in `.kind/`
  - Optionally install NGINX as ingress controller
  - Available `make` targets are prefixed with `kind-`

* CLI and logging framework
  - To help get you started with CLI subcommands, flags and environment variables
  - If you don't need subcommands, remove `example_command.go` and adjust `cli.App` settings in `main.go`

## TODO's after generating from this template

TIP: You can search for these tasks using `grep -n -r "TODO:" .`

1. `go.mod`: Adjust module name.
1. `.goreleaser.yml`: Adjust Docker image location in `dockers` and `docker_manifests` parameters.
1. `.gitignore`: Adjust binary file name.
1. `Dockerfile`: Adjust binary file name.
1. `Makefile.vars.mk`: Adjust project meta.
1. `.github/ISSUE_TEMPLATE/config.yml` (optional): Enable forwarding questions to GitHub Discussions or other page.
1. `.github/workflows/test.yml`: Update CodeClimate reporter ID (to be found in codeclimate.com Test coverage settings)
1. `docs/antora.yml`: Adjust project meta.
1. `docs/antora-playbook.yml`: Adjust project meta.
1. `docs/modules/pages/index.adoc`: Edit start page.
1. `docs/modules/nav.adoc`: Edit navigation.
1. `main.go`: Adjust variables.
1. Edit this README (including badges links)
1. Start hacking in `example_command.go`.

After completing a task, you can remove the comment in the files.

## Other repository settings

1. GitHub Settings
   - "Options > Wiki" (disable)
   - "Options > Allow auto-merge" (enable)
   - "Options > Automatically delete head branches" (enable)
   - "Collaborators & Teams > Add Teams and users to grant maintainer permissions
   - "Branches > Branch protection rules":
     - Branch name pattern: `master`
     - Require status check to pass before merging: `["lint"]` (you may need to push come commits first)
   - "Pages > Source": Branch `gh-pages`

1. GitHub Issue labels
   - "Issues > Labels > New Label" for the following labels with color suggestions:
     - `change` (`#D93F0B`)
     - `dependency` (`#ededed`)
     - `breaking` (`#FBCA04`)

1. CodeClimate Settings
   - "Repo Settings > GitHub > Pull request status updates" (install)
   - "Repo Settings > Test coverage > Enforce {Diff,Total} Coverage" (configure to your liking)

## Antora documentation

This template comes with an Antora documentation module to help you create Asciidoctor documentation.
By default, it is automatically published to GitHub Pages in `gh-pages` branch, however it can also be included in external Antora playbooks.

### Setup GitHub Pages

Once you generated a new repository using this template, the initial commit automatically runs a Job that creates the documentation in the `gh-pages` branch.
All you need to do is then to enable Pages in the settings.

The `gh-pages` branch is a parent-less commit that only contains the Antora-generated files.

However, if that's not the case or if you are setting up Antora in an existing repository, here's how you can achieve the same, but make sure to **commit or stash current changes first!**
```bash
current_branch=$(git rev-parse --abbrev-ref HEAD)
initial_commit=$(git rev-list --max-parents=0 HEAD | tail -n 1)
git switch --create gh-pages $initial_commit
git rm -r --cached .
git commit -m "Prepare gh-pages branch"
git push --set-upstream origin gh-pages
git switch -f $current_branch
```

And you're done!
GitHub automatically recognizes activity and sets up Pages if there's a `gh-pages` branch.

---

If you want to skip deployment to GitHub Pages you need to delete specific files and references:
`rm -f .github/workflows/docs.yml docs/package*.json docs/antora-playbook.yml docs/antora-build.mk`.
Also don't forget to delete the branch and disable Pages in the repository settings.

---

If you want to remove documentation completely simply run `rm -rf docs .github/workflows/docs.yml`.
