package gen_svc

import (
	"os"
	"os/exec"
	"path"
)

func genGit(serviceName, svcPath string) error {
	outData, err := os.ReadFile("templates/.gitignore.template")
	if err != nil {
		return err
	}

	filePath := path.Join(svcPath, ".gitignore")
	if err := os.WriteFile(filePath, outData, 0666); err != nil {
		return err
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = "../core-" + serviceName
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
