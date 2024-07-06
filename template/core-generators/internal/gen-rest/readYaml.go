package gen_rest

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYaml(serviceName string) (*Rest, error) {
	data, err := os.ReadFile(fmt.Sprintf("../{{.prjName}}-schema/%s/rest.yaml", serviceName))
	if err != nil {
		return nil, err
	}

	var f Rest

	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	return &f, nil
}
