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
		return errors.New("exec err: " + err.Error())
	}

	fmt.Println(stdoutBuf.String())
	fmt.Println(stderrBuf.String())

	return nil
}
