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
	"os"
	"regexp"
	"strings"

	"github.com/dtaniwaki/git-kustomize-diff/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	Base          string
	Target        string
	IncludeRegexp *regexp.Regexp
	ExcludeRegexp *regexp.Regexp
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

	log.Infof("Clone the git repo at %s for base", baseCommit)
	baseDirPath, err := ioutil.TempDir("", "git-kustomize-diff-base-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(baseDirPath)
	baseGitDir, err := currentGitDir.CloneAndCheckout(baseDirPath, baseCommit)
	if err != nil {
		return err
	}

	log.Infof("Clone the git repo at %s for target", baseCommit)
	targetDirPath, err := ioutil.TempDir("", "git-kustomize-diff-target-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(targetDirPath)
	targetGitDir, err := currentGitDir.CloneAndCheckout(targetDirPath, baseCommit)
	if err != nil {
		return err
	}
	log.Infof("Merge the commit at %s into the target repo", targetCommit)
	err = targetGitDir.Merge(targetCommit)
	if err != nil {
		return err
	}

	diffMap, err := Diff(baseGitDir.WorkDir.Dir, targetGitDir.WorkDir.Dir, DiffOpts{
		IncludeRegexp: opts.IncludeRegexp,
		ExcludeRegexp: opts.ExcludeRegexp,
	})
	if err != nil {
		return err
	}

	dirs := diffMap.Dirs()
	fmt.Printf("# Git Kustomize Diff\n\n")
	fmt.Println("| name | value |")
	fmt.Println("|-|-|")
	fmt.Printf("| dir | %s |\n", dirPath)
	fmt.Printf("| base | %s |\n", opts.Base)
	fmt.Printf("| target | %s |\n", opts.Target)
	fmt.Println("")

	fmt.Printf("## Target Kustomizations\n\n")
	if len(dirs) > 0 {
		fmt.Printf("```\n%s\n```\n", strings.Join(dirs, "\n"))
	} else {
		fmt.Println("N/A")
	}
	fmt.Println("")

	fmt.Printf("## Diff\n\n")
	lines := make([]string, 0)
	for idx, dir := range dirs {
		text := diffMap.Results[dir].AsMarkdown()
		if text != "" {
			lines[idx] = fmt.Sprintf("### %s:\n%s", dir, text)
		}
	}
	if len(lines) > 0 {
		fmt.Println(strings.Join(lines, "\n"))
	} else {
		fmt.Println("N/A")
	}
	fmt.Println("")

	return nil
}
