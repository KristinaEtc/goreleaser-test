package goreleaser

import (
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goreleaser/goreleaser/config"
	"github.com/pkg/errors"
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
	setReleaseDefaults(ctx)
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

func setReleaseDefaults(ctx *Context) error {
	if ctx.Config.Release.GitHub.Name != "" {
		return nil
	}
	repo, err := remoteRepo()
	if err != nil {
		return err
	}
	ctx.Config.Release.GitHub = repo
	return nil
}

// remoteRepo gets the repo name from the Git config.
func remoteRepo() (result config.Repo, err error) {
	if !IsRepo() {
		return result, errors.New("current folder is not a git repository")
	}
	out, err := gitRun("config", "--get", "remote.origin.url")
	if err != nil {
		return result, errors.Wrap(err, "repository doesn't have an `origin` remote")
	}
	return extractRepoFromURL(out), nil
}

// Run runs a git command and returns its output or errors
func gitRun(args ...string) (output string, err error) {
	var cmd = exec.Command("git", args...)
	bts, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(bts))
	}
	return string(bts), err
}

func IsRepo() bool {
	out, err := gitRun("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

func extractRepoFromURL(s string) config.Repo {
	for _, r := range []string{
		"git@github.com:",
		".git",
		"https://github.com/",
		"\n",
	} {
		s = strings.Replace(s, r, "", -1)
	}
	return toRepo(s)
}

func toRepo(s string) config.Repo {
	var ss = strings.Split(s, "/")
	return config.Repo{
		Owner: ss[0],
		Name:  ss[1],
	}
}

// AddArtifact adds a file to upload list
func (ctx *Context) AddArtifact(file string) {
	artifactsLock.Lock()
	defer artifactsLock.Unlock()
	file = strings.TrimPrefix(file, ctx.Config.Dist+string(filepath.Separator))
	ctx.Artifacts = append(ctx.Artifacts, file)
}
