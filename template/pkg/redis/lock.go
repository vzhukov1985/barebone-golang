package redis

import (
	"github.com/{{.orgName}}/{{.pkgRepoName}}/errors"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"time"
)

const defaultLockTimeout = time.Minute

func Lock(name string, timeout *time.Duration) (func(), error) {
	m := rs.NewMutex(name)

	t := defaultLockTimeout
	if timeout != nil {
		t = *timeout
	}
	checkTime := time.Now().Add(t)

	for {
		err := m.Lock()
		if err != nil {
			if timeout != nil && time.Now().After(checkTime) {
				log.Errorw("Failed to set mutex lock", "error", err)
				return nil, errors.LockError.WithInfo("lock timeout")
			}
			time.Sleep(50 * time.Millisecond)
		} else {
			break
		}
	}

	return func() {
		if ok, err := m.Unlock(); !ok || err != nil {
			log.Errorw("Failed to unlock redis mutex", "error", err)
		}
	}, nil
}
