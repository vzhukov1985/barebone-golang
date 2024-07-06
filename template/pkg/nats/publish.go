package nats

import (
	"github.com/goccy/go-json"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
)

func Publish(subj string, payload interface{}) {
	var data []byte
	if payload != nil {
		var err error
		data, err = json.Marshal(payload)
		if err != nil {
			log.Errorw("Failed to marshal nats message payload", "error", err)
			return
		}
	}

	err := Conn.Publish(subj, data)
	if err != nil {
		log.Errorw("Failed to send nats message", "error", err, "subj", subj)
	}
	log.Debugw("Nats message published", "subj", subj)
}
