package app

import (
	"github.com/{{.orgName}}/{{.pkgRepoName}}/env"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
)

func LogStart() {
	log.Infof("Starting %s service...", Name)
	log.Infof("Environment: %s", env.Current())
}

func LogStop() {
	log.Infof("Terminating %s service...", Name)
}
