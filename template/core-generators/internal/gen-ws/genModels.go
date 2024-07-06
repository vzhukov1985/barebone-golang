package gen_ws

import (
	gen_rest "core-generators/internal/gen-rest"
	"os"
)

func GenModels(api *Api) error {
	outPath := "../core-websocket/internal/models/ws-models"

	if err := os.RemoveAll(outPath); err != nil {
		return err
	}

	if err := os.MkdirAll(outPath, 0777); err != nil {
		return err
	}

	for _, v := range api.Models {
		if err := gen_rest.GenModel(outPath, "wsModels", v); err != nil {
			return err
		}
	}

	return nil
}
