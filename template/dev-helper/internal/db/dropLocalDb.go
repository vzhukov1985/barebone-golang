package db

import (
	"context"
	db "github.com/{{.orgName}}/{{.pkgRepoName}}/mongo"
)

func DropLocalDb() error {
	return db.Db.Drop(context.Background())
}
