package goreleaser

import (
	"bytes"
	"log"
	"text/template"

	"github.com/KristinaEtc/goreleaser/buildtarget"
	"github.com/goreleaser/goreleaser/config"
)

// from https://github.com/goreleaser/goreleaser

type nameData struct {
	Os          string
	Arch        string
	Arm         string
	Version     string
	Tag         string
	ProjectName string
	Binary      string
}

// ForBuild return the name for the given context, goos, goarch, goarm and
// build, using the build.Binary property instead of project_name.
func ForBuild(ctx *Context, build config.Build, target buildtarget.Target) (string, error) {
	return apply(
		nameData{
			Os:      replace(ctx.Config.Archive.Replacements, target.OS),
			Arch:    replace(ctx.Config.Archive.Replacements, target.Arch),
			Arm:     replace(ctx.Config.Archive.Replacements, target.Arm),
			Version: ctx.Version,
			Tag:     ctx.Git.CurrentTag,
			//Binary:      build.Binary,
			ProjectName: "goreleaser",
			Binary:      "goreleaser",
			//	ProjectName: build.Binary,
		},
		ctx.Config.Archive.NameTemplate,
	)
}

// ForName returns the name for the given context, goos, goarch and goarm.
func ForName(ctx *Context, target buildtarget.Target) (string, error) {
	return apply(
		nameData{
			Os:      replace(ctx.Config.Archive.Replacements, target.OS),
			Arch:    replace(ctx.Config.Archive.Replacements, target.Arch),
			Arm:     replace(ctx.Config.Archive.Replacements, target.Arm),
			Version: ctx.Version,
			Tag:     ctx.Git.CurrentTag,
			//Binary:      ctx.Config.ProjectName,
			//	ProjectName: ctx.Config.ProjectName,
			ProjectName: "goreleaser",
			Binary:      "goreleaser",
		},
		//ctx.Config.Archive.NameTemplate,
		NameTemplate,
	)
}

// ForChecksums returns the filename for the checksums file based on its
// template
func ForChecksums(ctx *Context) (string, error) {
	return apply(
		nameData{
			//	ProjectName: ctx.Config.ProjectName,
			Tag:     ctx.Git.CurrentTag,
			Version: ctx.Version,
		},
		ctx.Config.Checksum.NameTemplate,
	)
}

// ForTitle returns the release title based upon its template
func ForTitle(ctx *Context) (string, error) {
	return apply(
		nameData{
			ProjectName: ctx.Config.ProjectName,
			Tag:         ctx.Git.CurrentTag,
			Version:     ctx.Version,
		},
		ctx.Config.Release.NameTemplate,
	)
}

func apply(data nameData, templateStr string) (string, error) {
	//	log.Printf("data=%+v\n", data)
	//log.Println("template.New")
	//	log.Printf("templatestr=%s\n", templateStr)
	//os.Exit(1)
	var out bytes.Buffer
	t, err := template.New(data.ProjectName).Parse(templateStr)
	if err != nil {
		log.Println("template.New err=", err.Error())
		return "", err
	}
	err = t.Execute(&out, data)
	if err != nil {
		log.Println("Execute err=", err.Error())
	}
	//	log.Println("out.String(=", out.String())
	return out.String(), err
}

func replace(replacements map[string]string, original string) string {
	result := replacements[original]
	if result == "" {
		return original
	}
	return result
}
