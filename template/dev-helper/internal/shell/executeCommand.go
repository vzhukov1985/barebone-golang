package shell

import (
	"bufio"
	"context"
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"os"
	"os/exec"
)

type ExecuteCommandOpts struct {
	WorkDir        string
	EnvFile        string
	SendDataFunc   func(data string)
	ShowTerminated bool
}

func ExecuteCommand(opts ExecuteCommandOpts, cmd string, args ...string) {
	if err := func() error {
		ctx, _ := context.WithCancel(context.Background())
		shCmd := exec.CommandContext(ctx, cmd, args...)
		shCmd.Dir = opts.WorkDir

		if opts.EnvFile != "" {
			envs, err := readLines(opts.EnvFile)
			if err != nil {
				return err
			}

			shCmd.Env = append(os.Environ(), envs...)
		}

		stdout, err := shCmd.StdoutPipe()
		if err != nil {
			return err
		}
		shCmd.Stderr = shCmd.Stdout

		if err = shCmd.Start(); err != nil {
			return err
		}

		if opts.SendDataFunc != nil {
			s := bufio.NewScanner(stdout)
			for s.Scan() {
				opts.SendDataFunc(s.Text())
			}
		}

		if opts.ShowTerminated {
			opts.SendDataFunc("!!!!!!!!!!! Terminated !!!!!!!!!!!!!")
		}

		return nil
	}(); err != nil {
		log.Error(err)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
