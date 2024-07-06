package utils

import (
	"dev-helper/internal/services"
	"dev-helper/internal/shell"
	"fmt"
	"sync"
)

func PullAll() {
	svcs := append([]string{"remember-schema", "core-packages", "core-build-deploy"}, services.Slice...)
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

func PullService(serviceName string) {
	opts := shell.ExecuteCommandOpts{
		WorkDir:      fmt.Sprintf("../%s", serviceName),
		SendDataFunc: func(data string) { fmt.Println(serviceName + ": " + data) },
	}

	shell.ExecuteCommand(opts, "git", "pull")
}
