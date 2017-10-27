package goreleaser

import (
	"io/ioutil"

	"log"

	"github.com/goreleaser/goreleaser/config"
	yaml "gopkg.in/yaml.v2"
)

const (
	confFile              = ".goreleaser.yml"
	dirWithGitHubReleases = ".dist"
	secret                = "secret"
)

func Release(conf *config.Project) error {
	log.Printf("DEBUG conf: %+v\n", conf)
	if err := RunBuild(conf); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func ParseGitHubConf(confName string) (*config.Project, error) {
	b, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	var conf config.Project
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, err
}

func PushToGithub() error {
	return nil
}
