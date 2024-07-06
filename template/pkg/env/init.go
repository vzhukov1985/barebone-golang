package env

import (
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"slices"
)

var environment = Development

func init() {
	environment = Environment(GetString("ENVIRONMENT", "Development", false))
	if !slices.Contains(EnvironmentSlice, environment) {
		environment = Development
	}

	switch environment {
	case Development, Stage:
		log.SetLevel(log.DebugLevel)
	case Production:
		log.SetLevel(log.InfoLevel)
	}
}
