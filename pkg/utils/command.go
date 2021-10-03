/*
Copyright 2021 Daisuke Taniwaki.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type WithFormat interface {
	Format(fmt.State, rune)
}

type WithCause interface {
	Cause() error
}

type CommandError struct {
	InternalError error
	Stdout        string
	Stderr        string
}

func (ce *CommandError) Error() string {
	return fmt.Sprintf("%s\n\nstdout:\n%s\n\nstderr:\n%s", ce.InternalError, ce.Stdout, ce.Stderr)
}

func (ce *CommandError) Cause() error {
	return ce.InternalError.(WithCause).Cause()
}

func (ce *CommandError) Format(s fmt.State, verb rune) {
	ce.InternalError.(WithFormat).Format(s, verb)
}

func (ce *CommandError) ExitCode() *int {
	err, ok := ce.Cause().(*exec.ExitError)
	if !ok {
		return nil
	}
	code := err.ExitCode()
	return &code
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
			InternalError: errors.WithStack(err),
			Stdout:        stdout.String(),
			Stderr:        stderr.String(),
		}
	}
	return stdout.String(), stderr.String(), err
}
