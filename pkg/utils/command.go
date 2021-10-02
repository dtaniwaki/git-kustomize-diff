package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type CommandError struct {
	InternalError error
	Stdout        string
	Stderr        string
}

func (ce *CommandError) Error() string {
	return fmt.Sprintf("%s\n\nstdout:\n%s\n\nstderr:\n%s", ce.InternalError, ce.Stdout, ce.Stderr)
}

type WorkDir struct {
	Dir string
	Env map[string]string
}

func (wd *WorkDir) RunCommand(command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	cmd.Dir = wd.Dir
	env := make([]string, 0, len(os.Environ())+len(wd.Env)+1)
	env = append(env, os.Environ()...)
	for key, val := range wd.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, val))
	}
	cmd.Env = env
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		err = &CommandError{
			InternalError: err,
			Stdout:        stdout.String(),
			Stderr:        stderr.String(),
		}
	}
	return stdout.String(), stderr.String(), err
}
