package gen_rest

import (
	"core-generators/internal/utils"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type PkgImport struct {
	Name string
	Path string
}

func GenModels(serviceName string, schemas Schemas) error {
	filePath := fmt.Sprintf("../core-%s/internal/models/rest", serviceName)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(filePath); err != nil {
		return err
	}

	if len(schemas) <= 0 {
		return nil
	}
	if err := os.MkdirAll(path.Join(wd, filePath), 0777); err != nil {
		return err
	}

	for _, m := range schemas {
		if err := GenModel(filePath, "restModels", m); err != nil {
			return nil
		}
	}

	return nil
}

func GenModel(outPath string, packageName string, m Model) error {
	if m.Title == "" {
		return errors.New("model title not specified in schema")
	}
	m.Title = strings.ToUpper(m.Title[:1]) + m.Title[1:]

	var outData string

	outData += utils.GenModelHeader
	outData += "package " + packageName + "\n\n"
	oData, eData, pData := GenType(m.Title, m, true, true, true)

	outData += "import (\n"
	for _, p := range pData {
		outData += fmt.Sprintf("\t%s \"%s\"\n", p.Name, p.Path)
	}

	outData += ")\n"

	outData += oData

	outData += eData

	fileName := fmt.Sprintf("%s%s.go", strings.ToLower(m.Title[:1]), m.Title[1:])
	if err := os.WriteFile(path.Join(outPath, fileName), []byte(outData), 0666); err != nil {
		return err
	}

	utils.FormatFile(path.Join(outPath, fileName))

	return nil
}

func GenType(title string, m Model, isRequired bool, isRoot bool, isInternalModel bool) (string, string, []PkgImport) {
	imports := make([]PkgImport, 0)
	outData := ""

	if isRoot {
		outData += fmt.Sprintf("type %s ", title)
	}

	if !isRequired && m.Type != "object" && m.Type != "array" {
		outData += "*"
	}

	enumData := "\n"

	if m.Title == "" {
		m.Title = strings.ToUpper(title[:1]) + title[1:] + "Enum"
	}

	switch m.Type {
	case "string":
		if !isRoot && m.Enum != nil && len(m.Enum) > 0 {
			outData += m.Title
		} else {
			outData += "string"
		}
	case "boolean":
		if !isRoot && m.Enum != nil && len(m.Enum) > 0 {
			outData += m.Title
		} else {
			outData += "bool"
		}
	case "integer":
		if !isRoot && m.Enum != nil && len(m.Enum) > 0 {
			outData += m.Title
		} else {
			outData += "int"
		}
	case "number":
		if !isRoot && m.Enum != nil && len(m.Enum) > 0 {
			outData += m.Title
		} else {
			outData += "float64"
		}
	case "object":
		oData, eData, pData := genObj(m, isInternalModel)
		outData += oData
		enumData += eData
		imports = append(imports, pData...)
	case "array":
		oData, eData, pData := genArray(title, m, isInternalModel)
		outData += oData
		enumData += eData
		imports = append(imports, pData...)
	}

	if m.Ref != "" {
		var refModelsPrefix string
		modelName := utils.UpperCaseFirstLetter(filepath.Base(m.Ref))
		if strings.HasPrefix(m.Ref, "#") {
			if !isInternalModel {
				refModelsPrefix = "restModels"
			}
		} else {
			modelName = strings.TrimSuffix(modelName, filepath.Ext(modelName))
			subDir := filepath.Base(filepath.Dir(m.Ref))
			if subDir == "models" {
				refModelsPrefix = "commonModels"
			} else {
				if subDir != "." {
					if subDir == ".." {
						refModelsPrefix = "models"
					} else {
						refModelsPrefix = subDir + "Models"
					}
				}
			}
			if subDir != "." {
				pth := strings.ReplaceAll(filepath.Dir(m.Ref), "..", "")
				if !strings.HasPrefix(pth, "/models") {
					pth = "/models/" + pth
				}
				imports = append(imports, PkgImport{
					Name: refModelsPrefix,
					Path: "github.com/{{.orgName}}/{{.pkgRepoName}}" + pth,
				})
			}

		}

		if refModelsPrefix == "" {
			outData += modelName
		} else {
			outData += refModelsPrefix + "." + modelName
		}
	}

	if m.AnyOf != nil {
		outData += "interface{}"
	}

	enumData += genEnums(m, !isRoot)

	return outData, enumData, imports
}

func genEnums(m Model, genType bool) string {
	if m.Enum != nil {
		res := ""

		if genType {
			res += "type " + m.Title + " " + m.Type
		}

		res += "\n"
		res += "const (\n"
		for _, e := range m.Enum {
			var enumTitle string
			if m.Type == "number" {
				enumTitle = strings.ReplaceAll(e, ".", "p")
			} else {
				enumTitle = e
			}
			res += "\t" + m.Title + strings.ToUpper(enumTitle[:1]) + enumTitle[1:] + " " + m.Title + " = "
			if m.Type == "string" {
				res += fmt.Sprintf("\"%s\"", e) + "\n"
			} else {
				res += e + "\n"
			}
		}
		res += ")\n\n"

		res += fmt.Sprintf("var %sSlice = []%s {\n", m.Title, m.Title)
		for _, e := range m.Enum {
			if m.Type == "number" {
				e = strings.ReplaceAll(e, ".", "p")
			}
			res += "\t" + m.Title + strings.ToUpper(e[:1]) + e[1:] + ",\n"
		}
		res += "}\n\n"

		res += fmt.Sprintf("func (e %s) Validate() bool {\n", m.Title)
		res += fmt.Sprintf("\tfor _, v := range %sSlice {\n", m.Title)
		res += `		if e == v {
			return true
		}
	}
	return false
}

`
		return res
	}
	return ""
}

func genObj(m Model, isInternalModel bool) (string, string, []PkgImport) {
	res := "struct {\n"
	eData := ""
	pData := make([]PkgImport, 0)

	if m.Properties != nil {
		type mSt struct {
			Name  string
			Model Model
		}

		pArr := make([]mSt, 0)
		for name, p := range m.Properties {
			pArr = append(pArr, mSt{
				Name:  name,
				Model: p,
			})
		}

		sort.Slice(pArr, func(i, j int) bool {
			return pArr[i].Name < pArr[j].Name
		})

		for _, p := range pArr {
			isRequired := false
			if m.Required != nil {
				for _, r := range m.Required {
					if r == p.Name {
						isRequired = true
						break
					}
				}
			}
			omitEmpty := ""
			if !isRequired {
				omitEmpty = ",omitempty"
			}

			oD, eD, pD := GenType(p.Name, p.Model, isRequired, false, isInternalModel)
			eData += eD
			pData = append(pData, pD...)
			res += "\t" + strings.ToUpper(p.Name[:1]) + p.Name[1:] + " " + oD + " `json:\"" + p.Name + omitEmpty + "\" bson:\"" + toSnakeCase(p.Name) + omitEmpty + "\"`\n"
		}
	}

	res += "}"
	return res, eData, pData
}

func genArray(title string, m Model, isInternalModel bool) (string, string, []PkgImport) {
	if m.Items == nil {
		panic("no items specified for array in scheme")
	}

	oData, eData, pData := GenType(title, *m.Items, true, false, isInternalModel)

	return "[]" + oData, eData, pData
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
