package main

import (
	"log"
	"os"

	"github.com/KristinaEtc/goreleaser"
)

func main() {

	conf, err := goreleaser.ParseGitHubConf("../example/.goreleaser.yml")
	if err != nil {
		log.Println("ERR [ParseGitHubConf]: ", err.Error())
	}

	if err := goreleaser.Release(conf); err != nil {
		log.Println("ERR [Release]: ", err.Error())
		os.Exit(1)
	}
	log.Println("INFO Done successfully")
}
