package utils

import (
	"dev-helper/internal/db"
	"dev-helper/internal/shell"
	"fmt"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"os"
	"strings"
)

func RestoreCleanDb() {
	files, err := os.ReadDir("../core-build-deploy/BACKUPS/clean")
	if err != nil {
		log.Errorw("Failed to read clean backups dir", "error", err)
		return
	}

	if len(files) == 0 {
		log.Error("No clean backups found")
		return
	}

	var cleanFileName string
	for _, v := range files {
		if strings.HasSuffix(v.Name(), ".tar.gz") {
			cleanFileName = v.Name()
			break
		}
	}

	if cleanFileName == "" {
		log.Error("No clean backups found")
		return
	}

	opts := shell.ExecuteCommandOpts{
		WorkDir: "../core-build-deploy/BACKUPS",
		SendDataFunc: func(data string) {
			fmt.Println(data)
		},
		ShowTerminated: false,
	}

	shell.ExecuteCommand(opts, "cp", "./clean/"+cleanFileName, "clean_auto.tar.gz")

	log.Info("Dropping local db...")
	if err := db.DropLocalDb(); err != nil {
		log.Errorw("Failed to drop local db", "error", err)
		return
	}
	log.Info("Dropping local db successful")

	shell.ExecuteCommand(opts, "rm", "-rf", "db_local")
	shell.ExecuteCommand(opts, "tar", "-xf", "clean_auto.tar.gz")
	shell.ExecuteCommand(opts, "mongorestore", "--db=db_local", "--port=27017", "./db_local", "--gzip")
	shell.ExecuteCommand(opts, "rm", "-f", "clean_auto.tar.gz")
}
