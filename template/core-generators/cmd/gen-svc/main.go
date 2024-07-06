package main

import (
	gen_svc "core-generators/internal/gen-svc"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Укажи название сервиса в параметрах: make gen-svc имя_сервиса(без префикса core-) [доп.пар-ры]. 

Параметры:
	rest сгенерит дополнительно сервер рест
	worker сгенерит дополнительно костяк воркера

Примеры: "make gen-svc test rest worker" сгенерит костяк для нового сервиса core-test с рест сервером и воркером
         "make gen-svc test сгенерит только костяк нового сервиса core-test"`)

		return
	}

	serviceName := os.Args[1]

	isGenRest := false
	isGenWorker := false

	for _, a := range os.Args {
		switch a {
		case "rest":
			isGenRest = true
		case "worker":
			isGenWorker = true
		}
	}

	if err := gen_svc.GenService(serviceName, isGenRest, isGenWorker); err != nil {
		fmt.Println("Error:", err)
	}
}
