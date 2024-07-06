package shell

import "fmt"

func ExecuteBackendService(serviceName string) {
	ExecuteCommand(ExecuteCommandOpts{
		WorkDir:        fmt.Sprintf("../%s", serviceName),
		EnvFile:        fmt.Sprintf("../%s/.envs/.env.development.local", serviceName),
		SendDataFunc:   func(data string) { fmt.Println(serviceName + ": " + data) },
		ShowTerminated: true,
	}, "go", "run", "./cmd")
}
