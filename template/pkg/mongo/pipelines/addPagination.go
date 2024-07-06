package pipelines

import (
	"go.mongodb.org/mongo-driver/bson"
)

func AddPagination(pipeline *bson.A, skip, limit int) {
	stages := bson.A{}
	if skip > 0 {
		stages = append(stages, bson.D{{"$skip", skip}})
	}

	if limit > 0 {
		stages = append(stages, bson.D{{"$limit", limit}})
	}

	*pipeline = append(*pipeline, stages...)
}
