package mongo

import (
	"context"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Users *mongo.Collection
)

const (
	UsersCollection = "users"
)

var collectionNamesSlice = []string{
	UsersCollection,
}

func initCollections() {
	Users = Db.Collection(UsersCollection)

	createCollections()
	createIndexes()
}

func createCollections() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalw("Failed to create mongo db collections", "error", err)
		}
	}()

	existingCols, err := Db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return
	}

	for _, c := range collectionNamesSlice {
		isExist := false
		for _, exC := range existingCols {
			if c == exC {
				isExist = true
				break
			}
		}

		if !isExist {
			if err = Db.CreateCollection(context.Background(), c); err != nil {
				return
			}
		}
	}

	return
}

func createIndexes() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalw("Failed to create mongo db indexes", "error", err)
		}
	}()
}
