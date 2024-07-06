package prepare

import (
	"fmt"
	"golang/utils"
	"path/filepath"
)

func DevHelper(outDir, orgName string) error {
	wd := filepath.Join(outDir, "dev-helper")
	envs := []string{fmt.Sprintf("GOPRIVATE=github.com/%s", orgName)}
	if err := utils.ShellExec(wd, envs, "go", "mod", "init", "dev-helper"); err != nil {
		return err
	}

	if err := utils.ShellExec(wd, envs, "go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}
