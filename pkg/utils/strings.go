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
	"io/ioutil"
	"os/exec"
	"strings"
)

func Diff(text1, text2 string) (string, error) {
	tmpFile1, err := ioutil.TempFile("", "kustomize-diff-diff-")
	if err != nil {
		return "", err
	}
	defer tmpFile1.Close()
	_, err = tmpFile1.Write([]byte(text1))
	if err != nil {
		return "", err
	}

	tmpFile2, err := ioutil.TempFile("", "kustomize-diff-diff-")
	if err != nil {
		return "", err
	}
	defer tmpFile2.Close()
	_, err = tmpFile2.Write([]byte(text2))
	if err != nil {
		return "", err
	}

	stdout, _, err := (&WorkDir{}).RunCommand("diff", "-u", tmpFile1.Name(), tmpFile2.Name())
	if err != nil {
		_, ok := err.(*CommandError).InternalError.(*exec.ExitError)
		if !ok {
			return "", err
		}
	}

	lines := strings.Split(stdout, "\n")
	if len(lines) < 2 {
		// no diff
		return "", nil
	}
	return strings.Join(lines[2:], "\n"), nil
}
