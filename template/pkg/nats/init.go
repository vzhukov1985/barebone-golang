package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/env"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"time"
)

var Conn *nats.Conn

type Msg = nats.Msg
type Subscription = nats.Subscription

func init() {
	uri := env.GetString("NATS_URI", "", true)
	var err error
	Conn, err = nats.Connect(uri, nats.Timeout(time.Second*5))
	if err != nil {
		log.Fatalw("Failed to connect to NATS", "error", err, "uri", uri)
	}
	log.Infow("Successfully connected to NATS", "uri", uri)
}
