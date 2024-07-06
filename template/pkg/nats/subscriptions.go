package nats

import (
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
)

type Subscriptions []*Subscription

func (ss Subscriptions) UnsubscribeAll() {
	for _, v := range ss {
		if err := v.Unsubscribe(); err != nil {
			if v == nil {
				log.Errorw("Subscription is nil")
				continue
			}
			log.Errorw("Failed to unsubscribe from nats message", "error", err, "subj", v.Subject)
		}
	}
}
