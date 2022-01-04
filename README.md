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
  - Go build for amd64
  - Docker build for `latest` and `vx.y.z` tags
  - Push Docker image to GitHub's registry `ghcr.io`

## TODO's after generating from this template

TIP: You can search for these tasks using `grep -n -r "TODO:" .`

1. `go.mod`: Adjust module name.
1. `.goreleaser.yml`: Adjust Docker image location in `dockers` parameter.
1. `.gitignore`: Adjust binary file name.
1. `Dockerfile`: Adjust binary file name.
1. `.github/ISSUE_TEMPLATE/config.yml` (optional): Enable forwarding questions to GitHub Discussions or other page.

After doing the tasks, you can remove the comments in the files.
