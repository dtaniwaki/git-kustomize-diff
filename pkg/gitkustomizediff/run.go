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
	"io/ioutil"
	"os"
	"regexp"

	"github.com/dtaniwaki/git-kustomize-diff/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	Base                    string
	Target                  string
	IncludeRegexp           *regexp.Regexp
	ExcludeRegexp           *regexp.Regexp
	KustomizePath           string
	KustomizeLoadRestrictor string
	GitPath                 string
	Debug                   bool
	AllowDirty              bool
}

type RunResult struct {
	BaseCommit   string
	TargetCommit string
	DiffMap      *DiffMap
}

func Run(dirPath string, opts RunOpts) (*RunResult, error) {
	log.Info("Start run")
	currentGitDir := utils.NewGitDir(dirPath, opts.GitPath)
	baseCommitish := opts.Base
	if baseCommitish == "" {
		baseCommitish = "origin/main"
	}
	baseCommit, err := currentGitDir.CommitHash(baseCommitish)
	if err != nil {
		return nil, err
	}
	targetCommitish := opts.Target
	if targetCommitish == "" {
		targetCommitish, err = currentGitDir.CurrentBranch()
		if err != nil {
			return nil, err
		}
	}
	targetCommit, err := currentGitDir.CommitHash(targetCommitish)
	if err != nil {
		return nil, err
	}

	dirtyPatch := ""
	if opts.AllowDirty {
		log.Infof("Generate a dirty patch from %s", targetCommit)
		diff, err := currentGitDir.Diff(targetCommit)
		if err != nil {
			return nil, err
		}
		dirtyPatch = diff
	}

	log.Infof("Clone the git repo at %s for base", baseCommit)
	baseDirPath, err := ioutil.TempDir("", "git-kustomize-diff-base-")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if opts.Debug {
		log.Infof("Base repo path: %s", baseDirPath)
	} else {
		defer os.RemoveAll(baseDirPath)
	}
	baseGitDir, err := currentGitDir.CloneAndCheckout(baseDirPath, baseCommit)
	if err != nil {
		return nil, err
	}

	log.Infof("Clone the git repo at %s for target", baseCommit)
	targetDirPath, err := ioutil.TempDir("", "git-kustomize-diff-target-")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if opts.Debug {
		log.Infof("Target repo path: %s", targetDirPath)
	} else {
		defer os.RemoveAll(targetDirPath)
	}
	targetGitDir, err := currentGitDir.CloneAndCheckout(targetDirPath, baseCommit)
	if err != nil {
		return nil, err
	}
	log.Infof("Merge the commit at %s into the target repo", targetCommit)
	err = targetGitDir.Merge(targetCommit)
	if err != nil {
		return nil, err
	}
	if dirtyPatch != "" {
		log.Infof("Apply the dirty patch")
		err = targetGitDir.Apply(dirtyPatch)
		if err != nil {
			return nil, err
		}
	}

	diffMap, err := Diff(baseGitDir.WorkDir.Dir, targetGitDir.WorkDir.Dir, DiffOpts{
		IncludeRegexp:           opts.IncludeRegexp,
		ExcludeRegexp:           opts.ExcludeRegexp,
		KustomizePath:           opts.KustomizePath,
		KustomizeLoadRestrictor: opts.KustomizeLoadRestrictor,
	})
	if err != nil {
		return nil, err
	}

	return &RunResult{
		BaseCommit:   baseCommit,
		TargetCommit: targetCommit,
		DiffMap:      diffMap,
	}, nil
}
