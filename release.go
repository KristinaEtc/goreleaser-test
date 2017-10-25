package goreleaser

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

// RunRelease the pipe
func RunRelease(ctx *Context) error {
	c, err := NewGitHub(ctx)
	if err != nil {
		return err
	}
	return doRun(ctx, c)
}

func doRun(ctx *Context, c Client) error {
	/*if !ctx.Publish {
		return pipeline.Skip("--skip-publish is set")
	}*/
	ctx.Git.CurrentTag = "v1.8.0"

	log.Println("tag", ctx.Git.CurrentTag,
		"repo", ctx.Config.Release.GitHub.String(),
		" creating or updating release")

	body, err := describeBody(ctx)
	if err != nil {
		return err
	}

	log.Println(ctx.Config.Release.GitHub.Name)
	releaseID, err := c.CreateRelease(ctx, body.String())
	if err != nil {
		return err
	}
	var g errgroup.Group

	for _, artifact := range ctx.Artifacts {

		err = upload(ctx, c, releaseID, artifact)
		if err != nil {
			log.Println(" upload err=", err.Error())
		}

	}
	return g.Wait()
}

func upload(ctx *Context, c Client, releaseID int, artifact string) error {
	var path = filepath.Join(ctx.Config.Dist, artifact)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	_, name := filepath.Split(path)
	log.Println("file", file.Name(), " name", name, " uploading to release")
	return c.Upload(ctx, releaseID, name, file)
}
