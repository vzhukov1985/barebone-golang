package mongo

import (
	"context"
	"fmt"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func CreateTransaction(f func(sessionContext SessionContext) error) (err error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := Db.Client().StartSession()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
			if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
				log.Errorw("Failed to abort db transaction", "error", abortErr)
				return
			}

		}
		session.EndSession(context.Background())
	}()

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(txnOpts); err != nil {
			return err
		}

		if err = f(sessionContext); err != nil {
			return err
		}

		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			log.Errorw("Failed to abort db transaction", "error", abortErr)
			err = abortErr
			return
		}
		return
	}
	return
}
