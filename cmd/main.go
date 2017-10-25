package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/KristinaEtc/goreleaser"
	"github.com/goreleaser/goreleaser/config"
	yaml "gopkg.in/yaml.v2"
)

const (
	confFile              = ".goreleaser.yml"
	dirWithGitHubReleases = ".dist"
	secret                = "secret"
)

func PushToGithub() error {
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

func Release() error {
	conf, err := ParseGitHubConf(confFile)
	if err != nil {
		return err
	}
	log.Printf("conf: %+v\n", conf)
	if err := goreleaser.RunBuild(conf); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func main() {
	if err := Release(); err != nil {
		log.Println("ERR [parseGitHubConf config] ", err.Error())
		os.Exit(1)
	}
	log.Println("Done successfully")
}
