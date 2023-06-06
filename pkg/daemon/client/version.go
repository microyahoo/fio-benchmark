package client

import "github.com/microyahoo/fio-benchmark/pkg/util/exec"

func FioVersion(executor exec.Executor) (string, error) {
	args := []string{"--version"}
	output, err := executor.ExecuteCommandWithTimeout(FioCommandsTimeout, FioTool, args...)
	if err != nil {
		return "", err
	}
	return output, nil
}
