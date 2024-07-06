package utils

import (
	"dev-helper/internal/services"
	"dev-helper/internal/shell"
	"fmt"
	"sync"
)

func PushAll() {
	var wg sync.WaitGroup
	for _, v := range services.Slice {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			PushService(s)
		}(v)
	}
	wg.Wait()
}

func PushService(serviceName string) {
	opts := shell.ExecuteCommandOpts{
		WorkDir:      fmt.Sprintf("../%s", serviceName),
		SendDataFunc: func(data string) { fmt.Println(serviceName + ": " + data) },
	}

	shell.ExecuteCommand(opts, "git", "pull")
	shell.ExecuteCommand(opts, "git", "add", ".")
	shell.ExecuteCommand(opts, "git", "commit", "-m", "upd")
	shell.ExecuteCommand(opts, "git", "push")
}
