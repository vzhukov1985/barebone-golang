package gen_svc

import (
	"core-generators/internal/utils"
	"fmt"
	"os"
)

func genMain(serviceName, filePath string, isRest, isWorker bool) error {
	var outData string

	outData += fmt.Sprintf(`package main

import (
	internal_rest "core-%s/internal/api/rest"
	"core-%s/internal/worker"

	"github.com/{{.orgName}}/{{.pkgRepoName}}/app"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/rest"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app.LogStart()

`, serviceName, serviceName)

	if isWorker {
		outData += `worker.Start()

`
	}

	if isRest {
		outData += `restController := internal_rest.CreateController()
	restServer := rest.CreateServer(restController)
	restServer.Start()

`
	}

	outData += `	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

<-done
app.LogStop()
`
	if isWorker {
		outData += `worker.Stop()
`
	}

	if isRest {
		outData += `restServer.Stop()
`
	}

	outData += "}"

	if err := os.WriteFile(filePath, []byte(outData), 0666); err != nil {
		return err
	}

	utils.FormatFile(filePath)

	return nil
}
