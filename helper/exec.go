package helper

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func Exec(bashCommand string, args []string) error {
	cmd := exec.Command(bashCommand, args...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		return errors.New(err.Error())
	}

	fmt.Println("stdout", stdoutBuf.String())
	fmt.Println("stderr", stderrBuf.String())

	return nil
}
