package nats

import (
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"strings"
)

func Subscribe[T any](subj string, f func(msg *Msg, subjParts []string, payload T), queue ...string) *Subscription {
	var sub *Subscription
	var payload T
	var err error

	if len(queue) > 0 && queue[0] != "" {
		sub, err = Conn.QueueSubscribe(subj, queue[0], subscribeInnerFunc(subj, payload, f))
	} else {
		sub, err = Conn.Subscribe(subj, subscribeInnerFunc(subj, payload, f))
	}
	if err != nil {
		log.Fatalw("Failed to subscribe to nats message", "error", err, "subj", subj)
	}

	return sub
}

func SubscribeM[T any](subjs []string, f func(msg *Msg, subjParts []string, payload T), queue ...string) Subscriptions {
	var subs Subscriptions
	var payload T

	if len(queue) > 0 && queue[0] != "" {
		for _, v := range subjs {
			sub, err := Conn.QueueSubscribe(v, queue[0], subscribeInnerFunc(v, payload, f))
			if err != nil {
				log.Fatalw("Failed to subscribe to nats message", "error", err, "subj", v)
			}
			subs = append(subs, sub)
		}
	} else {
		for _, v := range subjs {
			sub, err := Conn.Subscribe(v, subscribeInnerFunc(v, payload, f))
			if err != nil {
				log.Fatalw("Failed to subscribe to nats message", "error", err, "subj", v)
			}
			subs = append(subs, sub)
		}
	}

	return subs
}

func subscribeInnerFunc[T any](subj string, payload T, f func(msg *Msg, subjParts []string, payload T)) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		if msg.Data != nil && len(msg.Data) > 0 {
			err := json.Unmarshal(msg.Data, &payload)
			if err != nil {
				log.Errorw("Failed to unmarshal payload", "error", err, "subj", subj)
			}
		}

		subjParts := strings.Split(msg.Subject, ".")

		f(msg, subjParts, payload)
	}
}
