package app

import "github.com/{{.orgName}}/{{.pkgRepoName}}/env"

var (
	Name = ""
)

func init() {
	Name = env.GetString("SERVICE_NAME", "", true)
}
