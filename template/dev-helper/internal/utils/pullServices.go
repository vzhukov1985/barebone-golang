package utils

import (
	"dev-helper/internal/services"
	"sync"
)

func PullServices() {
	svcs := append([]string{"core-packages", "core-build-deploy"}, services.Slice...)
	var wg sync.WaitGroup
	for _, v := range svcs {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			PullService(s)
		}(v)
	}
	wg.Wait()
}
