package goreleaser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	//"github.com/KristinaEtc/goreleaser/name"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubClient struct {
	client *github.Client
}

func getToken() (string, error) {
	token, err := ioutil.ReadFile(".secret")
	if err != nil {
		log.Println(err)
		return "", err
	}
	s := string(token)
	sz := len(string(token))
	return s[:sz-1], nil
}

// NewGitHub returns a github client implementation
func NewGitHub(ctx *Context) (Client, error) {
	token, err := getToken()
	if err != nil {
		log.Println(err.Error())
	}
	ctx.Token = token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ctx.Token},
	)
	t, err := ts.Token()
	if err != nil {
		log.Println("token err=", err)
	} else {
		log.Printf("tok=[%+v]\n", t)
	}

	c := oauth2.NewClient(ctx.Context, ts)
	fmt.Println("done")
	//os.Exit(1)
	client := github.NewClient(c)
	if ctx.Config.GitHubURLs.API != "" {
		api, err := url.Parse(ctx.Config.GitHubURLs.API)
		if err != nil {
			return &githubClient{}, err
		}
		upload, err := url.Parse(ctx.Config.GitHubURLs.Upload)
		if err != nil {
			return &githubClient{}, err
		}
		client.BaseURL = api
		client.UploadURL = upload
	}

	return &githubClient{client}, nil
}

func (c *githubClient) CreateFile(
	ctx *Context,
	content bytes.Buffer,
	path string,
) (err error) {
	options := &github.RepositoryContentFileOptions{
		Committer: &github.CommitAuthor{
			Name:  github.String(ctx.Config.Brew.CommitAuthor.Name),
			Email: github.String(ctx.Config.Brew.CommitAuthor.Email),
		},
		Content: content.Bytes(),
		Message: github.String(
			ctx.Config.ProjectName + " version " + ctx.Git.CurrentTag,
		),
	}

	file, _, res, err := c.client.Repositories.GetContents(
		ctx,
		ctx.Config.Brew.GitHub.Owner,
		ctx.Config.Brew.GitHub.Name,
		path,
		&github.RepositoryContentGetOptions{},
	)
	if err != nil && res.StatusCode == 404 {
		_, _, err = c.client.Repositories.CreateFile(
			ctx,
			ctx.Config.Brew.GitHub.Owner,
			ctx.Config.Brew.GitHub.Name,
			path,
			options,
		)
		return
	}
	options.SHA = file.SHA
	_, _, err = c.client.Repositories.UpdateFile(
		ctx,
		ctx.Config.Brew.GitHub.Owner,
		ctx.Config.Brew.GitHub.Name,
		path,
		options,
	)
	return
}

func (c *githubClient) CreateRelease(ctx *Context, body string) (releaseID int, err error) {
	var release *github.RepositoryRelease
	releaseTitle, err := ForTitle(ctx)
	if err != nil {
		return 0, err
	}
	var data = &github.RepositoryRelease{
		Name:       github.String(releaseTitle),
		TagName:    github.String(ctx.Git.CurrentTag),
		Body:       github.String(body),
		Draft:      github.Bool(ctx.Config.Release.Draft),
		Prerelease: github.Bool(ctx.Config.Release.Prerelease),
	}
	release, _, err = c.client.Repositories.GetReleaseByTag(
		ctx,
		ctx.Config.Release.GitHub.Owner,
		ctx.Config.Release.GitHub.Name,
		ctx.Git.CurrentTag,
	)
	if err != nil {
		release, _, err = c.client.Repositories.CreateRelease(
			ctx,
			ctx.Config.Release.GitHub.Owner,
			ctx.Config.Release.GitHub.Name,
			data,
		)
	} else {
		release, _, err = c.client.Repositories.EditRelease(
			ctx,
			ctx.Config.Release.GitHub.Owner,
			ctx.Config.Release.GitHub.Name,
			release.GetID(),
			data,
		)
	}
	log.Println("url", release.GetHTMLURL(), "release updated")
	return release.GetID(), err
}

func (c *githubClient) Upload(
	ctx *Context,
	releaseID int,
	name string,
	file *os.File,
) (err error) {
	_, _, err = c.client.Repositories.UploadReleaseAsset(
		ctx,
		ctx.Config.Release.GitHub.Owner,
		ctx.Config.Release.GitHub.Name,
		releaseID,
		&github.UploadOptions{
			Name: name,
		},
		file,
	)
	return
}
