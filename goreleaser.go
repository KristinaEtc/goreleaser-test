package goreleaser

import (
	"io/ioutil"

	"log"

	"github.com/goreleaser/goreleaser/config"
	yaml "gopkg.in/yaml.v2"
)

const (
	confFile = ".goreleaser.yml"
	secret   = "secret"
)

var dirWithGitHubReleases = "dist/"

func Release(conf *config.Project) error {
	log.Printf("DEBUG conf: %+v\n", conf)
	cnt, err := RunBuild(conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := Archive(cnt); err != nil {
		log.Fatalln(err.Error())
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
