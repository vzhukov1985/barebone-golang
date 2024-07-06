package main

import (
	"dev-helper/internal/utils"
	"dev-helper/internal/worker"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "mod-build-all":
			utils.ModBuildAll()
			return
		case "pull-services":
			utils.PullServices()
			return
		case "pull-all":
			utils.PullAll()
			return
		case "push-all":
			utils.PushAll()
			return
		case "run-all":
			worker.Start()

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			<-done
			worker.Stop()
			return
		case "find-id":
			if len(os.Args) > 2 {
				utils.FindId(os.Args[2])
			} else {
				fmt.Println("Укажи id: make find-id {id}")
			}
			return
		case "restore-clean-db":
			utils.RestoreCleanDb()
			return
		}
	}

	fmt.Println("No command specified")
	fmt.Println("Available commands:")
	fmt.Println("mod-upd-all")
	fmt.Println("pull-services")
	fmt.Println("pull-all")
	fmt.Println("push-all")
	fmt.Println("run-all")
}
