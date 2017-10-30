package goreleaser

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goreleaser/goreleaser/config"
)

var artifactsLock sync.Mutex

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

func setDefaultValues(ctx *Context) {
	log.Println("DEBUG set default values")
	ctx.Config.Archive = config.Archive{
		Format: "tar.gz",
		FormatOverrides: []config.FormatOverride{
			{
				Goos:   "windows",
				Format: "zip",
			},
		},
	}
}

// AddArtifact adds a file to upload list
func (ctx *Context) AddArtifact(file string) {
	artifactsLock.Lock()
	defer artifactsLock.Unlock()
	file = strings.TrimPrefix(file, ctx.Config.Dist+string(filepath.Separator))
	ctx.Artifacts = append(ctx.Artifacts, file)
}
