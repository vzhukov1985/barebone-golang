package gen_rest

import (
	"core-generators/internal/utils"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sort"
	"strings"
)

func GenErrors() error {
	schemaErrorsPath := "../{{.prjName}}-schema/models/error/ErrorCode.yaml"

	pkgErrorsPath := "../pkg/errors/errors.go"

	os.Remove(pkgErrorsPath)

	data, err := os.ReadFile(schemaErrorsPath)
	if err != nil {
		return err
	}

	var m Model

	if err := yaml.Unmarshal(data, &m); err != nil {
		return err
	}

	var outData string

	outData += utils.GenModelHeader

	outData += "package errors\n\nvar (\n"

	type errorInfo struct {
		Name         string
		Desc         string
		IsBadRequest bool
	}

	errorInfos := make([]errorInfo, 0)

	for eName, eDesc := range m.ErrorDescs {
		isBadRequest := false

		if strings.HasPrefix(eDesc, "+") {
			isBadRequest = true
			eDesc = eDesc[1:]
		}

		errorInfos = append(errorInfos, errorInfo{
			Name:         eName,
			Desc:         eDesc,
			IsBadRequest: isBadRequest,
		})
	}

	sort.Slice(errorInfos, func(i, j int) bool {
		return errorInfos[i].Name < errorInfos[j].Name
	})

	for _, e := range errorInfos {
		outData += fmt.Sprintf("\t%s = New(\"%s\", \"%s\", %t)\n", utils.UpperCaseFirstLetter(e.Name), e.Name, e.Desc, e.IsBadRequest)

	}

	outData += ")\n\nvar Slice = []*SError{\n"

	for _, e := range errorInfos {
		outData += "\t" + utils.UpperCaseFirstLetter(e.Name) + ",\n"
	}

	outData += "}"

	if err := os.WriteFile(pkgErrorsPath, []byte(outData), 0666); err != nil {
		return err
	}

	utils.FormatFile(pkgErrorsPath)

	return nil
}
