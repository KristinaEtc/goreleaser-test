package goreleaser

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/config"
	"github.com/goreleaser/goreleaser/context"
	"github.com/pkg/errors"
)

const (
	//NameTemplate        = "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}"
	//	ReleaseNameTemplate = "{{.Tag}}"

	// NameTemplate default name_template for the archive.
	NameTemplate = "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}"

	// ReleaseNameTemplate is the default name for the release.
	ReleaseNameTemplate = "{{.Tag}}"

	// SnapshotNameTemplate represents the default format for snapshot release names.
	SnapshotNameTemplate = "SNAPSHOT-{{ .Commit }}"

	// ChecksumNameTemplate is the default name_template for the checksum file.
	ChecksumNameTemplate = "{{ .ProjectName }}_{{ .Version }}_checksums.txt"
)

func SetDefault(ctx *context.Context) error {
	ctx.Config.Dist = "dist"

	if ctx.Config.ProjectName == "" {
		ctx.Config.ProjectName = ctx.Config.Release.GitHub.Name
	}

	setBuildDefaults(ctx)

	if ctx.Config.Brew.Install == "" {
		var installs []string
		for _, build := range ctx.Config.Builds {
			if !isBrewBuild(build) {
				continue
			}
			installs = append(
				installs,
				fmt.Sprintf(`bin.install "%s"`, build.Binary),
			)
		}
		ctx.Config.Brew.Install = strings.Join(installs, "\n")
	}

	if ctx.Config.Brew.CommitAuthor.Name == "" {
		ctx.Config.Brew.CommitAuthor.Name = "goreleaserbot"
	}
	if ctx.Config.Brew.CommitAuthor.Email == "" {
		ctx.Config.Brew.CommitAuthor.Email = "goreleaser@carlosbecker.com"
	}

	err := setArchiveDefaults(ctx)
	log.WithField("config", ctx.Config).Debug("defaults set")
	setReleaseDefaults(ctx)
	return err
}

func isBrewBuild(build config.Build) bool {
	for _, ignore := range build.Ignore {
		if ignore.Goos == "darwin" && ignore.Goarch == "amd64" {
			return false
		}
	}
	return contains(build.Goos, "darwin") && contains(build.Goarch, "amd64")
}

func contains(ss []string, s string) bool {
	for _, zs := range ss {
		if zs == s {
			return true
		}
	}
	return false
}

func setBuildDefaults(ctx *context.Context) {
	for i, build := range ctx.Config.Builds {
		ctx.Config.Builds[i] = buildWithDefaults(ctx, build)
	}
	if len(ctx.Config.Builds) == 0 {
		ctx.Config.Builds = []config.Build{
			buildWithDefaults(ctx, ctx.Config.SingleBuild),
		}
	}
}

func buildWithDefaults(ctx *context.Context, build config.Build) config.Build {
	if build.Binary == "" {
		build.Binary = ctx.Config.Release.GitHub.Name
	}

	ctx.Config.Dist = "dist"
	if ctx.Config.Release.NameTemplate == "" {
		ctx.Config.Release.NameTemplate = ReleaseNameTemplate
	}

	if len(build.Goos) == 0 {
		build.Goos = []string{"linux", "darwin"}
	}
	if len(build.Goarch) == 0 {
		build.Goarch = []string{"amd64", "386"}
	}

	if build.Ldflags == "" {
		build.Ldflags = "-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}"
	}
	return build
}

func setArchiveDefaults(ctx *context.Context) error {
	if ctx.Config.Archive.NameTemplate == "" {
		ctx.Config.Archive.NameTemplate = NameTemplate
	}
	if ctx.Config.Archive.Format == "" {
		ctx.Config.Archive.Format = "tar.gz"
	}
	if len(ctx.Config.Archive.Files) == 0 {
		ctx.Config.Archive.Files = []string{
			"licence*",
			"LICENCE*",
			"license*",
			"LICENSE*",
			"readme*",
			"README*",
			"changelog*",
			"CHANGELOG*",
		}
	}
	return nil
}

func setReleaseDefaults(ctx *context.Context) error {
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
	if !gitIsRepo() {
		return result, errors.New("current folder is not a git repository")
	}
	out, err := gitRun("config", "--get", "remote.origin.url")
	if err != nil {
		return result, errors.Wrap(err, "repository doesn't have an `origin` remote")
	}
	return extractRepoFromURL(out), nil
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

// IsRepo returns true if current folder is a git repository
func gitIsRepo() bool {
	out, err := gitRun("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// gitRun runs a git command and returns its output or errors
func gitRun(args ...string) (output string, err error) {
	var cmd = exec.Command("git", args...)
	bts, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(bts))
	}
	return string(bts), err
}

// Clean the output
func Clean(output string, err error) (string, error) {
	return strings.Replace(strings.Split(output, "\n")[0], "'", "", -1), err
}
