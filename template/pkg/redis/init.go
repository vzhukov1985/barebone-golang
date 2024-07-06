package redis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/env"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
)

var client *redis.Client
var rs *redsync.Redsync

func init() {
	address := env.GetString("REDIS_URI", "", true)
	db := env.GetInt("REDIS_DB", 0, true)

	client = redis.NewClient(&redis.Options{
		Addr: address,
		DB:   db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalw("Failed to connect to redis", "error", err, "uri", address, "db", db)
	}

	pool := goredis.NewPool(client)
	rs = redsync.New(pool)
}
