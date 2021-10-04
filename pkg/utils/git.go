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
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/yookoala/realpath"
)

type GitDir struct {
	GitPath string
	WorkDir WorkDir
}

func NewGitDir(dirPath, gitPath string) *GitDir {
	return &GitDir{
		GitPath: gitPath,
		WorkDir: WorkDir{
			Dir: dirPath,
		},
	}
}

func (gd *GitDir) RunGitCommand(args ...string) (string, string, error) {
	gitPath := gd.GitPath
	if gitPath == "" {
		gitPath = "git"
	}
	return gd.WorkDir.RunCommand(gitPath, args...)
}

func (gd *GitDir) CommitHash(target string) (string, error) {
	stdout, _, err := gd.RunGitCommand("rev-parse", "-q", "--short", target)
	if err != nil {
		return "", err
	}
	return strings.Trim(stdout, "\n"), nil
}

func (gd *GitDir) Diff(target string) (string, error) {
	stdout, _, err := gd.RunGitCommand("diff", target)
	if err != nil {
		return "", err
	}
	return stdout, nil
}

func (gd *GitDir) CurrentBranch() (string, error) {
	stdout, _, err := gd.RunGitCommand("branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strings.Trim(stdout, "\n"), nil
}

func (gd *GitDir) Clone(dstDirPath string) (*GitDir, error) {
	rootDir, err := gd.GetRootDir()
	if err != nil {
		return nil, err
	}
	_, _, err = gd.RunGitCommand("clone", rootDir, dstDirPath)
	if err != nil {
		return nil, err
	}
	absPath, err := realpath.Realpath(gd.WorkDir.Dir)
	if err != nil {
		return nil, err
	}
	relPath, err := filepath.Rel(rootDir, absPath)
	if err != nil {
		return nil, err
	}
	return &GitDir{
		GitPath: gd.GitPath,
		WorkDir: WorkDir{Dir: filepath.Join(dstDirPath, relPath)},
	}, nil
}

func (gd *GitDir) GetRootDir() (string, error) {
	// `git rev-parse --show-toplevel` returns a real path.
	baseDirPath, _, err := gd.RunGitCommand("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.Trim(baseDirPath, "\n"), nil
}

func (gd *GitDir) CopyConfig(targetGitDir *GitDir) error {
	baseDirPath, err := gd.GetRootDir()
	if err != nil {
		return err
	}
	targetDirPath, err := targetGitDir.GetRootDir()
	if err != nil {
		return err
	}
	if baseDirPath == targetDirPath {
		return nil
	}
	src, err := os.Open(filepath.Join(baseDirPath, ".git", "config"))
	if err != nil {
		return errors.WithStack(err)
	}
	defer src.Close()
	dst, err := os.Create(filepath.Join(targetDirPath, ".git", "config"))
	if err != nil {
		return errors.WithStack(err)
	}
	defer dst.Close()

	// Manual copy as io.Copy copy_file_range has some problem in some situation.
	st, err := src.Stat()
	if err != nil {
		return errors.WithStack(err)
	}
	bs := make([]byte, st.Size())
	_, err = src.Read(bs)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = dst.Write(bs)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (gd *GitDir) Fetch() error {
	_, _, err := gd.RunGitCommand("fetch", "--all")
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) Checkout(target string) error {
	_, _, err := gd.RunGitCommand("checkout", target)
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) Merge(target string) error {
	_, _, err := gd.RunGitCommand("merge", "--no-ff", target)
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) Apply(patch string) error {
	tmpFile, err := ioutil.TempFile("", "git-kustomize-diff-apply-")
	if err != nil {
		return errors.WithStack(err)
	}
	defer (func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	})()
	_, err = tmpFile.Write([]byte(patch))
	if err != nil {
		return errors.WithStack(err)
	}
	_, _, err = gd.RunGitCommand("apply", tmpFile.Name())
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) SetUser() error {
	email := "anonymous@example.com"
	name := "anonymous"
	_, _, err := gd.RunGitCommand("config", "user.email", email)
	if err != nil {
		return err
	}
	_, _, err = gd.RunGitCommand("config", "user.name", name)
	if err != nil {
		return err
	}
	return nil
}

func (gd *GitDir) CloneAndCheckout(dirPath, commit string) (*GitDir, error) {
	gitDir, err := gd.Clone(dirPath)
	if err != nil {
		return nil, err
	}
	err = gd.CopyConfig(gitDir)
	if err != nil {
		return nil, err
	}
	err = gitDir.SetUser()
	if err != nil {
		return nil, err
	}
	err = gitDir.Fetch()
	if err != nil {
		return nil, err
	}
	err = gitDir.Checkout(commit)
	if err != nil {
		return nil, err
	}
	return gitDir, nil
}
