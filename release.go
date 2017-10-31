package goreleaser

import (
	"context"
	"log"
	"os"
	"path/filepath"
	//"github.com/KristinaEtc/goreleaser/client"
)

// RunRelease the pipe
func RunRelease(ctx *Context) error {
	c, err := NewGitHub(ctx)
	if err != nil {
		return err
	}
	ctx.Context = context.Background()
	return doRun(ctx, c)
}

func doRun(ctx *Context, c Client) error {
	/*	if !ctx.Publish {
		return pipeline.Skip("--skip-publish is set")
	}*/
	log.Println("tag", ctx.Git.CurrentTag,
		"repo", ctx.Config.Release.GitHub.String(), "creating or updating release")
	body, err := describeBody(ctx)
	if err != nil {
		return err
	}
	releaseID, err := c.CreateRelease(ctx, body.String())
	if err != nil {
		return err
	}

	for _, artifact := range ctx.Artifacts {
		//sem <- true
		artifact := artifact
		err = upload(ctx, c, releaseID, artifact)
	}
	return err
}

func upload(ctx *Context, c Client, releaseID int, artifact string) error {
	var path = filepath.Join(ctx.Config.Dist, artifact)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	_, name := filepath.Split(path)
	log.Println("file", file.Name(), "name", name, "uploading to release")
	return c.Upload(ctx, releaseID, name, file)
}
