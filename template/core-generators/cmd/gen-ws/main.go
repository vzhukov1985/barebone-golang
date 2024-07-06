package main

import (
	gen_ws "core-generators/internal/gen-ws"
	"fmt"
)

func main() {
	api, err := gen_ws.ReadYaml()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err := gen_ws.GenModels(api); err != nil {
		fmt.Println("GenModels Error: ", err)
		return
	}

	if err := gen_ws.GenHandlers(api); err != nil {
		fmt.Println("GenHandlers Error: ", err)
		return
	}

	fmt.Println("Websocket структура сгенерирована без ошибок")
}
