package gen_svc

import (
	"fmt"
	"os"
	"path"
)

func GenService(serviceName string, isRest, isWorker bool) error {
	svcPath := fmt.Sprintf("../core-%s", serviceName)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = os.Stat(path.Join(wd, svcPath))
	if err == nil {
		fmt.Println("Папка с сервисом под этим именем уже существует")
		return nil
	}

	if err := os.MkdirAll(path.Join(wd, svcPath, "cmd"), 0777); err != nil {
		return err
	}

	if err := genMain(serviceName, path.Join(wd, svcPath, "cmd", "main.go"), isRest, isWorker); err != nil {
		return err
	}

	if isWorker {
		if err := genWorker(serviceName, path.Join(wd, svcPath)); err != nil {
			return err
		}
	}

	if isRest {
		if err := genRest(serviceName, path.Join(wd, svcPath)); err != nil {
			return err
		}
	}

	if err := genMakefile(serviceName, path.Join(wd, svcPath)); err != nil {
		return err
	}

	if err := genGit(serviceName, path.Join(wd, svcPath)); err != nil {
		return err
	}

	if err := genEnvs(serviceName, path.Join(wd, svcPath), isRest); err != nil {
		return err
	}

	if err := genMod(serviceName); err != nil {
		return err
	}

	return nil
}
