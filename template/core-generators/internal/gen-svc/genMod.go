package gen_svc

import (
	"os/exec"
)

func genMod(serviceName string) error {
	cmd := exec.Command("go", "mod", "init", "core-"+serviceName)
	cmd.Dir = "../core-" + serviceName
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
