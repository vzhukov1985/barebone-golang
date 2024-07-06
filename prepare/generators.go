package prepare

import (
	"fmt"
	"golang/utils"
	"path/filepath"
)

func Generators(outDir, orgName string) error {
	wd := filepath.Join(outDir, "core-generators")
	envs := []string{fmt.Sprintf("GOPRIVATE=github.com/%s", orgName)}
	if err := utils.ShellExec(wd, envs, "go", "mod", "init", "core-generators"); err != nil {
		return err
	}

	if err := utils.ShellExec(wd, envs, "go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}
