package nats

import (
	errors2 "errors"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/app"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/errors"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"strings"
)

func SubscribeWithResponse[P any, R any](subj string, f func(msg *Msg, subjParts []string, payload P) (R, error)) *Subscription {
	procFunc := func(payload P) func(msg *nats.Msg) {
		return func(msg *nats.Msg) {
			respInnerPayload, err := func() ([]byte, error) {
				if msg.Data != nil && len(msg.Data) > 0 {
					err := json.Unmarshal(msg.Data, &payload)
					if err != nil {
						log.Errorw("Failed to unmarshal payload", "error", err, "subj", subj)
						return nil, err
					}
				}

				subjParts := strings.Split(msg.Subject, ".")

				resp, err := f(msg, subjParts, payload)
				if err != nil {
					return nil, err
				}
				var payloadBytes []byte
				payloadBytes, err = json.Marshal(resp)
				if err != nil {
					return nil, err
				}
				return payloadBytes, nil
			}()

			var sErr *errors.SError
			if err != nil {
				var e *errors.SError
				switch {
				case errors2.As(err, &e):
					sErr = e
				default:
					sErr = errors.Undefined.WithInfo(err.Error())
				}
			}

			p := json.RawMessage(respInnerPayload)
			msgData, err := json.Marshal(ResponsePayload{
				Payload: &p,
				Error:   sErr,
			})
			if err != nil {
				log.Errorw("Failed to marshall nats reply message", "error", err)
				return
			}
			if err = msg.Respond(msgData); err != nil {
				log.Errorw("Failed to send reply message", "error", err)
				return
			}
		}
	}

	var payload P
	sub, errOuter := Conn.QueueSubscribe(subj, app.Name, procFunc(payload))
	if errOuter != nil {
		log.Fatalw("Failed to subscribe with response to nats message", "error", errOuter, "subj", subj)
		return nil
	}

	return sub
}
