# go-bootstrap
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
After completing a task, you can remove the comment in the files.

## Other repository settings

1. GitHub Settings
   - "Options > Wiki" (disable)
   - "Options > Allow auto-merge" (enable)
   - "Options > Automatically delete head branches" (enable)
   - "Branches > Default branch" (change to `master`) (VSHN's default)
   - "Branches > Branch protection rules":
     - Branch name pattern: `master`
     - Require status check to pass before merging: `["lint"]` (you may need to push come commits first)
   -
1. CodeClimate Settings
   - "Repo Settings > GitHub > Pull request status updates" (install)
   - "Repo Settings > Test coverage > Enforce {Diff,Total} Coverage" (configure to your liking)

## Antora documentation

This template comes with an Antora documentation module to help you create sophisticated Asciidoctor documentation.
By default, it is automatically published to GitHub Pages in `gh-pages` branch, however it can also be included in external Antora playbooks.

### Setup GitHub Pages

1. Prepare `gh-pages` branch (commit current changes first!)
```
initial_commit=$(git rev-list --max-parents=0 HEAD)
git switch --track --create gh-pages $initial_commit
#rm -rf *
```

If you want to skip publish to GitHub Pages, simply remove `.github/workflows/docs.yml`.

If you want to remove documentation completely follow these steps:

1. Remove `docs`
1. Remove `.github/workflows/docs.yml`
1. Remove line `include docs/Makefile` in `Makefile`
