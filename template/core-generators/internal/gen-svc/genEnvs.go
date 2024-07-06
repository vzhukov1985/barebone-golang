package gen_svc

import (
	"os"
	"path"
)

func genEnvs(serviceName, svcPath string, isRest bool) error {
	err := os.MkdirAll(path.Join(svcPath, ".envs"), 0777)
	if err != nil {
		return err
	}

	data := "SERVICE_NAME=" + serviceName + "\n"
	data += "ENVIRONMENT=development\n\n"

	if isRest {
		data += "REST_ADDRESS="
	}

	if err := os.WriteFile(path.Join(svcPath, ".envs", ".env.template"), []byte(data), 0666); err != nil {
		return err
	}

	return nil
}
