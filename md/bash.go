package md

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// RunCMD runs arbitary bash command with arguments.
// If erOutIsError is set to true, everything that goes to
// error std out will be treated as an error.
func RunCMD(command string, erOutIsError bool) (string, error) {
	bash_cmd_list := strings.Split(command, " ")
	var cmd *exec.Cmd
	if len(bash_cmd_list) == 0 {
		return "", errors.New("command not found")
	} else if len(bash_cmd_list) == 1 {
		cmd = exec.Command(bash_cmd_list[0])
	} else {
		cmd = exec.Command(bash_cmd_list[0], bash_cmd_list[1:]...)
	}
	var stderr bytes.Buffer
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		s := ""
		if stderr.String() != "" {
			s = ": " + stderr.String()
		}
		return "", fmt.Errorf("error: %v%s", err, s)

	} else if erOutIsError && stderr.String() != "" {
		return "", fmt.Errorf("%s", stderr.String())
	}
	return out.String(), nil
}
