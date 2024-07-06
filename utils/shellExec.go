package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func ShellExec(workDir string, envs []string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), envs...)
	out, err := cmd.Output()
	if len(out) > 0 {
		fmt.Println(string(out))
	}
	if err != nil {
		var execError *exec.ExitError
		if ok := errors.As(err, &execError); ok {
			if len(execError.Stderr) > 0 {
				fmt.Println(string(execError.Stderr))
			}
			fmt.Println("Exit code:", execError.ExitCode())
		}
		return err
	}
	return nil
}
