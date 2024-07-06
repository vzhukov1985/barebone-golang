package main

import (
	gen_rest "core-generators/internal/gen-rest"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Укажи название сервиса в параметрах: make gen-rest имя_сервиса(без префикса core-). Пример: \"make gen-rest auth\" сгенерит рест для сервиса core-auth")
		return
	}

	serviceName := os.Args[1]

	f, err := gen_rest.ReadYaml(serviceName)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err := gen_rest.GenModels(serviceName, f.Components["schemas"]); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err := gen_rest.GenHandlers(serviceName, f.Paths); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("REST структура сгенерирована без ошибок для сервиса", serviceName)

	if err := gen_rest.GenCommonModels(); err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Общие модели сгенерированы без ошибок")

	if err := gen_rest.GenErrors(); err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Коды ошибок сгенерированы без ошибок")
}
