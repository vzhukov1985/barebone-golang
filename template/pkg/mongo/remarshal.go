package mongo

import "go.mongodb.org/mongo-driver/bson"

func Remarshal(in interface{}, out interface{}) error {
	bytes, err := bson.Marshal(in)
	if err != nil {
		return err
	}
	return bson.Unmarshal(bytes, out)
}
