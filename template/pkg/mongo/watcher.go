package mongo

import (
	"context"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WatchDoc struct {
	DocumentId      string
	OperationType   OperationType
	UpdatedFields   bson.M
	RemovedFields   bson.A
	TruncatedArrays bson.A
	NewDocument     bson.M
	OldDocument     bson.M
}

type Watcher struct {
	coll  *mongo.Collection
	funcs []func(d *WatchDoc)
}

type OperationType string

const (
	OpInsert OperationType = "insert"
	OpUpdate OperationType = "update"
	OpDelete OperationType = "delete"
)

func CreateWatcher(collection *mongo.Collection, p ...func(d *WatchDoc)) *Watcher {
	return &Watcher{coll: collection, funcs: p}
}

func (w *Watcher) Start() {
	w.startDb(Db)
}

func (w *Watcher) StartDb(db *mongo.Database) {
	w.startDb(db)
}

func (w *Watcher) startDb(db *mongo.Database) {
	go func() {
		res := db.RunCommand(context.Background(), bson.D{
			{"collMod", w.coll.Name()},
			{"changeStreamPreAndPostImages", bson.D{{"enabled", true}}},
		})
		if res.Err() != nil {
			log.Error(res.Err())
		}

		opts := options.ChangeStream()
		opts.SetFullDocumentBeforeChange(options.WhenAvailable)
		//opts.SetFullDocument(options.WhenAvailable)

		cs, err := w.coll.Watch(context.Background(), bson.D{}, opts)
		if err != nil {
			log.Fatalw("Failed to start watching collection", "error", err, "collection", w.coll.Name())
		}
		log.Infof("Started watching %s collection", w.coll.Name())

		for cs.Next(context.Background()) {
			var dbDoc struct {
				OperationType string `bson:"operationType"`
				DocumentKey   struct {
					Id string `bson:"_id"`
				} `bson:"documentKey"`
				UpdateDescription struct {
					UpdatedFields   bson.M `bson:"updatedFields"`
					RemovedFields   bson.A `bson:"removedFields"`
					TruncatedArrays bson.A `bson:"truncatedArrays"`
				} `bson:"updateDescription"`
				FullDocument             bson.M `bson:"fullDocument"`
				FullDocumentBeforeChange bson.M `bson:"fullDocumentBeforeChange"`
			}

			//var dd bson.M
			//cs.Decode(&dd)
			d := &WatchDoc{}
			err = cs.Decode(&dbDoc)
			if err != nil {
				log.Errorw("error decoding change stream document", "error", err, "collection", w.coll.Name())
			}
			d.DocumentId = dbDoc.DocumentKey.Id
			d.OperationType = OperationType(dbDoc.OperationType)
			d.UpdatedFields = dbDoc.UpdateDescription.UpdatedFields
			d.RemovedFields = dbDoc.UpdateDescription.RemovedFields
			d.TruncatedArrays = dbDoc.UpdateDescription.TruncatedArrays
			d.NewDocument = dbDoc.FullDocument
			d.OldDocument = dbDoc.FullDocumentBeforeChange

			for _, f := range w.funcs {
				f(d)
			}
		}
	}()
}
