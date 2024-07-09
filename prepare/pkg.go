package prepare

import (
	"fmt"
	"golang/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Pkg(outDir, prjName, orgName, pkgRepoName string, pkgOrder []string) error {
	if err := renameInsideTemplates(outDir, orgName, pkgRepoName, prjName); err != nil {
		return err
	}

	if err := renameSchemaDir(outDir, prjName); err != nil {
		return err
	}

	if err := pkgInitGit(outDir, orgName, pkgRepoName); err != nil {
		return err
	}

	if err := createModules(outDir, orgName, pkgRepoName, pkgOrder); err != nil {
		return err
	}

	return nil
}

func renameInsideTemplates(outDir string, orgName string, pkgRepoName, prjName string) error {
	if err := filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fDataBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			fOutDataStr := string(fDataBytes)

			fOutDataStr = strings.ReplaceAll(fOutDataStr, "{{.orgName}}", orgName)
			fOutDataStr = strings.ReplaceAll(fOutDataStr, "{{.pkgRepoName}}", pkgRepoName)
			fOutDataStr = strings.ReplaceAll(fOutDataStr, "{{.prjName}}", prjName)
			if err := os.WriteFile(path, []byte(fOutDataStr), 0644); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func renameSchemaDir(outDir, prjName string) error {
	if err := os.Rename(path.Join(outDir, "schema"), path.Join(outDir, prjName+"-schema")); err != nil {
		return err
	}

	return nil
}

func pkgInitGit(outDir, orgName, pkgRepoName string) error {
	pkgDir := filepath.Join(outDir, "pkg")

	envs := make([]string, 0)
	if err := utils.ShellExec(pkgDir, envs, "git", "init"); err != nil {
		return err
	}

	if err := utils.ShellExec(pkgDir, envs, "git", "add", "."); err != nil {
		return err
	}

	if err := utils.ShellExec(pkgDir, envs, "git", "commit", "-m", "initial commit"); err != nil {
		return err
	}

	if err := utils.ShellExec(pkgDir, envs, "git", "branch", "-M", "master"); err != nil {
		return err
	}

	if err := utils.ShellExec(pkgDir, envs, "git", "remote", "add", "origin", fmt.Sprintf("https://github.com/%s/%s.git", orgName, pkgRepoName)); err != nil {
		return err
	}

	if err := utils.ShellExec(pkgDir, envs, "git", "push", "-u", "origin", "master"); err != nil {
		return err
	}

	return nil
}

func createModules(outDir, orgName, pkgRepoName string, pkgOrder []string) error {
	pkgDir := filepath.Join(outDir, "pkg")
	envs := []string{fmt.Sprintf("GOPRIVATE=github.com/%s", orgName)}
	for _, v := range pkgOrder {
		wd := filepath.Join(pkgDir, v)
		if err := utils.ShellExec(wd, envs, "go", "mod", "init", fmt.Sprintf("github.com/%s/%s/%s", orgName, pkgRepoName, v)); err != nil {
			return err
		}

		if err := utils.ShellExec(wd, envs, "go", "mod", "tidy"); err != nil {
			return err
		}

		if err := utils.ShellExec(wd, []string{}, "git", "add", "."); err != nil {
			return err
		}

		if err := utils.ShellExec(wd, []string{}, "git", "commit", "-m", "initial version"); err != nil {
			return err
		}

		if err := utils.ShellExec(wd, []string{}, "git", "push"); err != nil {
			return err
		}

		if err := utils.ShellExec(wd, []string{}, "git", "tag", fmt.Sprintf("%s/v1.0.0", v)); err != nil {
			return err
		}

		if err := utils.ShellExec(wd, []string{}, "git", "push", "origin", fmt.Sprintf("%s/v1.0.0", v)); err != nil {
			return err
		}
	}
	return nil
}
