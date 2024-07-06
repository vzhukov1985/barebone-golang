package utils

import (
	"context"
	"fmt"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	db "github.com/{{.orgName}}/{{.pkgRepoName}}/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindId(id string) {
	var err error
	defer func() {
		if err != nil {
			log.Error("Error: ", err)
		}
	}()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	collections, err := db.Db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return
	}

	for _, collection := range collections {
		col := db.Db.Collection(collection)
		cursor, err := col.Find(context.TODO(), bson.M{})
		if err != nil {
			return
		}

		for cursor.Next(context.TODO()) {
			var results bson.M
			err := cursor.Decode(&results)
			if err != nil {
				return
			}

			processField(results, objId, collection, results["_id"].(primitive.ObjectID), "")
		}
	}
}

func processField(f interface{}, objId primitive.ObjectID, collection string, docId primitive.ObjectID, field string) {
	_, ok := f.(primitive.ObjectID)
	if ok {
		if f == objId {
			fmt.Printf("Collection: %s , Document ID: %s , Field: %s\n", collection, docId.Hex(), field)
		}
	}

	array, ok := f.(bson.A)
	for i, arrElem := range array {
		processField(arrElem, objId, collection, docId, fmt.Sprintf("%s.[%v]", field, i))
	}

	obj, ok := f.(bson.M)
	for k, v := range obj {
		var newField string
		if field == "" {
			newField = k
		} else {
			newField = fmt.Sprintf("%s.%s", field, k)
		}
		processField(v, objId, collection, docId, newField)
	}
}
