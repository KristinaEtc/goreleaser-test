package goreleaser

import (
	"context"

	"github.com/goreleaser/goreleaser/config"
)

// Context carries along some data through the pipes
type Context struct {
	context.Context
	Config       config.Project
	Token        string
	Git          GitInfo
	Binaries     map[string]map[string][]Binary
	Artifacts    []string
	ReleaseNotes string
	Version      string
	Validate     bool
	Publish      bool
}

// Binary with pretty name and path
type Binary struct {
	Name, Path string
}

type GitInfo struct {
	CurrentTag string
	Commit     string
}
