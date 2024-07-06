package mongo

import (
	"context"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/env"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionContext mongo.SessionContext

var Db *mongo.Database

func init() {
	uri := env.GetString("MONGO_DB_CONNECTION", "", true)
	dbName := env.GetString("MONGO_DB_NAME", "", true)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalw("Failed to connect to mongo db", "error", err, "uri", uri, "dbName", dbName)
		return
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalw("Failed to ping check mongo db", "error", err, "uri", uri, "dbName", dbName)
		return
	}

	Db = client.Database(dbName)
	initCollections()
}
