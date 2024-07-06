package pipelines

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CountAggregationDocuments(c *mongo.Collection, pipeline bson.A) (int, error) {
	cPipeline := bson.A{}
	cPipeline = append(cPipeline, pipeline...)
	cPipeline = append(cPipeline, bson.D{{"$count", "total"}})

	cur, err := c.Aggregate(context.Background(), cPipeline)
	if err != nil {
		return 0, err
	}

	var total int

	if cur.Next(context.Background()) {
		var d struct {
			Total int `bson:"total"`
		}

		if err := cur.Decode(&d); err != nil {
			return 0, err
		}

		total = d.Total
	}
	return total, nil
}
