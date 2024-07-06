package gen_rest

import (
	"core-generators/internal/utils"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GenCommonModels() error {
	schemaModelsPath := "../{{.prjName}}-schema/models"

	pkgModelsPath := "../pkg/models"

	if err := cleanModelsDir(pkgModelsPath); err != nil {
		return err
	}

	if err := processModelsSchema(schemaModelsPath, pkgModelsPath); err != nil {
		return err
	}

	return nil
}

func cleanModelsDir(path string) error {
	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			strData := string(data)
			if strings.HasPrefix(strData, utils.GenModelHeader) {
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && info.Name() != "models" {
			isEmpty, err := utils.IsDirEmpty(path)
			if err != nil {
				return err
			}
			if isEmpty {
				err := os.Remove(path)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func processModelsSchema(schemaModelsPath, pkgModelsPath string) error {
	if err := filepath.Walk(schemaModelsPath, func(entryPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() != "error" && info.Name() != "models" {
				err := os.MkdirAll(path.Join(pkgModelsPath, info.Name()), 0777)
				if err != nil {
					return err
				}
			}
		} else {
			if err := processModelFile(entryPath, pkgModelsPath); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return err
	}

	return nil
}

func processModelFile(filePath string, pkgModelsPath string) error {
	var subFolder string
	var pkgName string
	if path.Dir(filePath) == "../{{.prjName}}-schema/models" {
		pkgName = "models"
	} else {
		pkgName = utils.LowerCaseFirstLetter(filepath.Base(filepath.Dir(filePath)))
		subFolder = strings.ReplaceAll(path.Dir(filePath), "../{{.prjName}}-schema/models/", "")
	}

	if filepath.Ext(filePath) != ".yaml" {
		return nil
	}

	if subFolder == "error" {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var m Model

	if err := yaml.Unmarshal(data, &m); err != nil {
		return err
	}

	return GenModel(path.Join(pkgModelsPath, subFolder), pkgName, m)
}
