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
	"io"
	"os"
	"path/filepath"
	"strings"
)

type GitDir struct {
	workDir WorkDir
}

func (gd *GitDir) CommitHash(target string) (string, error) {
	stdout, _, err := gd.workDir.RunCommand("git", "rev-parse", "-q", "--short", target)
	if err != nil {
		return "", err
	}
	return strings.Trim(stdout, "\n"), nil
}

func (gd *GitDir) Clone(dstDirPath string) error {
	_, _, err := gd.workDir.RunCommand("git", "clone", gd.workDir.Dir, dstDirPath)
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) CopyConfig(dstDirPath string) error {
	src, err := os.Open(filepath.Join(gd.workDir.Dir, ".git", "config"))
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Open(filepath.Join(dstDirPath, ".git", "config"))
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(src, dst)
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) Fetch() error {
	_, _, err := gd.workDir.RunCommand("git", "fetch", "--all")
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) Checkout(target string) error {
	_, _, err := gd.workDir.RunCommand("git", "checkout", target)
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) Merge(target string) error {
	_, _, err := gd.workDir.RunCommand("git", "merge", "--no-ff", target)
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) SetUser() error {
	email := "anonymous@example.com"
	name := "anonymous"
	_, _, err := gd.workDir.RunCommand("git", "config", "user.email", email)
	if err != nil {
		return err
	}
	_, _, err = gd.workDir.RunCommand("git", "config", "user.name", name)
	if err != nil {
		return err
	}
	return nil
}
