package utils

import (
	"dev-helper/internal/services"
	"dev-helper/internal/shell"
	"fmt"
	"sync"
)

func ModBuildAll() {
	var wg sync.WaitGroup
	for _, v := range services.Slice {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			ModBuild(s)
		}(v)
	}
	wg.Wait()
}

func ModBuild(serviceName string) {
	opts := shell.ExecuteCommandOpts{
		WorkDir:      fmt.Sprintf("../%s", serviceName),
		SendDataFunc: func(data string) { fmt.Println(serviceName + ": " + data) },
	}

	shell.ExecuteCommand(opts, "make", "mod-build")
}
