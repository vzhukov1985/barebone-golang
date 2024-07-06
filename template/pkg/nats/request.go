package nats

import (
	"github.com/{{.orgName}}/{{.pkgRepoName}}/errors"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"time"
)

type ResponsePayload struct {
	Payload *json.RawMessage
	Error   *errors.SError
}

const StdRequestTimeout = 10 * time.Second

func Request[T any](subj string, reqPayload any, timeout ...time.Duration) (T, error) {
	var resp T
	var data []byte
	if reqPayload != nil {
		var err error
		data, err = json.Marshal(reqPayload)
		if err != nil {
			log.Errorw("Failed to marshal nats message payload", "error", err)
			return resp, errors.ParseRequest.WithError(err)
		}
	}

	to := StdRequestTimeout
	if len(timeout) > 0 {
		to = timeout[0]
	}

	msg, err := Conn.Request(subj, data, to)
	if err != nil {
		return resp, err
	}

	var reply ResponsePayload

	if err = json.Unmarshal(msg.Data, &reply); err != nil {
		return resp, err
	}

	if reply.Error != nil {
		return resp, reply.Error
	}
	if reply.Payload != nil {
		if err = json.Unmarshal(*reply.Payload, &resp); err != nil {
			return resp, err
		}
	}

	return resp, nil
}
