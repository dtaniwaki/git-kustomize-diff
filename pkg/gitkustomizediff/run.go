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

package gitkustomizediff

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dtaniwaki/git-kustomize-diff/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	Base   string
	Target string
}

func Run(dirPath string, opts RunOpts) error {
	log.Info("Start run")
	currentGitDir := &utils.GitDir{
		WorkDir: utils.WorkDir{
			Dir: dirPath,
		},
	}
	baseCommitish := opts.Base
	if baseCommitish == "" {
		baseCommitish = "origin/main"
	}
	baseCommit, err := currentGitDir.CommitHash(baseCommitish)
	if err != nil {
		return err
	}
	targetCommitish := opts.Target
	if targetCommitish == "" {
		targetCommitish, err = currentGitDir.CurrentBranch()
		if err != nil {
			return err
		}
	}
	targetCommit, err := currentGitDir.CommitHash(targetCommitish)
	if err != nil {
		return err
	}

	log.Info("Clone the git repo for base")
	baseDirPath, err := ioutil.TempDir("", "git-kustomize-diff-base-")
	if err != nil {
		return err
	}
	baseGitDir, err := currentGitDir.Clone(baseDirPath)
	if err != nil {
		return err
	}
	err = currentGitDir.CopyConfig(baseGitDir)
	if err != nil {
		return err
	}
	err = baseGitDir.Fetch()
	if err != nil {
		return err
	}
	err = baseGitDir.Checkout(baseCommit)
	if err != nil {
		return err
	}

	log.Info("Clone the git repo for target")
	targetDirPath, err := ioutil.TempDir("", "git-kustomize-diff-target-")
	if err != nil {
		return err
	}
	targetGitDir, err := currentGitDir.Clone(targetDirPath)
	if err != nil {
		return err
	}
	err = currentGitDir.CopyConfig(targetGitDir)
	if err != nil {
		return err
	}
	err = targetGitDir.Fetch()
	if err != nil {
		return err
	}
	err = targetGitDir.Checkout(baseCommit)
	if err != nil {
		return err
	}
	err = targetGitDir.Merge(targetCommit)
	if err != nil {
		return err
	}

	diffMap, err := Diff(baseGitDir.WorkDir.Dir, targetGitDir.WorkDir.Dir)
	if err != nil {
		return err
	}

	dirs := diffMap.Dirs()
	fmt.Printf("# Git Kustomize Diff")
	fmt.Printf("## Target Kustomizations\n\n```\n%s\n```\n\n", strings.Join(dirs, "\n"))

	fmt.Printf("## Diff\n")
	lines := make([]string, len(dirs))
	for idx, path := range dirs {
		text := diffMap.Results[path].AsMarkdown()
		if text != "" {
			lines[idx] = fmt.Sprintf("## %s:\n%s", path, text)
		}
	}
	fmt.Println(strings.Join(lines, "\n"))

	return nil
}
