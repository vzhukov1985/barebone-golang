package gen_svc

import (
	"os"
	"path"
	"strings"
)

func genMakefile(serviceName, svcPath string) error {
	tmpl, err := os.ReadFile("templates/Makefile.template")
	outDataStr := strings.Replace(string(tmpl), "{{.serviceName}}", serviceName, -1)

	if err != nil {
		return err
	}

	filePath := path.Join(svcPath, "Makefile")
	if err := os.WriteFile(filePath, []byte(outDataStr), 0666); err != nil {
		return err
	}

	return nil
}
