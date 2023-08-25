//go:build tools
// +build tools

// Package tools is a place to put any tooling dependencies as imports.
// Go modules will be forced to download and install them.
package tools

import (
	_ "github.com/crossplane/crossplane-tools/cmd/angryjet"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)
