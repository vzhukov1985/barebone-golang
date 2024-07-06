package utils

import (
	"fmt"
	"os/exec"
)

func FormatFile(filePath string) {
	arg0 := "-w"
	arg1 := filePath

	cmd := exec.Command("gofmt", arg0, arg1)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cmd = exec.Command("goimports", arg0, arg1)
	_, err = cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
