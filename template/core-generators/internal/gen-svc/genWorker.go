package gen_svc

import (
	"core-generators/internal/utils"
	"os"
	"path"
)

func genWorker(serviceName, svcPath string) error {
	outData := `package worker

func Start() {
}

func Stop() {

}`

	workerPath := path.Join(svcPath, "internal", "worker")
	if err := os.MkdirAll(workerPath, 0777); err != nil {
		return err
	}

	filePath := path.Join(workerPath, "worker.go")
	if err := os.WriteFile(filePath, []byte(outData), 0666); err != nil {
		return err
	}

	utils.FormatFile(filePath)

	return nil
}
