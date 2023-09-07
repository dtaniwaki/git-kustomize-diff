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
	"path/filepath"
	"regexp"

	"github.com/dtaniwaki/git-kustomize-diff/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type DiffOpts struct {
	IncludeRegexp           *regexp.Regexp
	ExcludeRegexp           *regexp.Regexp
	KustomizePath           string
	KustomizeLoadRestrictor string
}

func Diff(baseDirPath, targetDirPath string, opts DiffOpts) (*DiffMap, error) {
	log.Info("Start diff")
	listOpts := utils.ListKustomizeDirsOpts{
		IncludeRegexp: opts.IncludeRegexp,
		ExcludeRegexp: opts.ExcludeRegexp,
	}
	baseKDirs, err := utils.ListKustomizeDirs(baseDirPath, listOpts)
	if err != nil {
		return nil, err
	}
	log.Debugf("base dirs: %+v", baseKDirs)
	targetKDirs, err := utils.ListKustomizeDirs(targetDirPath, listOpts)
	if err != nil {
		return nil, err
	}
	log.Debugf("target dirs: %+v", targetKDirs)
	kDirs := map[string]struct{}{}
	for _, kDir := range append(baseKDirs, targetKDirs...) {
		kDirs[kDir] = struct{}{}
	}
	diffMap := NewDiffMap()
	for kDir := range kDirs {
		baseKDirPath := filepath.Join(baseDirPath, kDir)
		if !utils.KustomizationExists(baseKDirPath) {
			err := utils.MakeKustomizeDir(baseKDirPath)
			if err != nil {
				diffMap.Results[kDir] = &DiffError{err}
				continue
			}
		}
		targetKDirPath := filepath.Join(targetDirPath, kDir)
		if !utils.KustomizationExists(targetKDirPath) {
			err := utils.MakeKustomizeDir(targetKDirPath)
			if err != nil {
				diffMap.Results[kDir] = &DiffError{err}
				continue
			}
		}
		baseYaml, err := Build(baseKDirPath, BuildOpts{opts.KustomizePath, opts.KustomizeLoadRestrictor})
		if err != nil {
			diffMap.Results[kDir] = &DiffError{err}
			continue
		}
		targetYaml, err := Build(targetKDirPath, BuildOpts{opts.KustomizePath, opts.KustomizeLoadRestrictor})
		if err != nil {
			diffMap.Results[kDir] = &DiffError{err}
			continue
		}

		content, err := utils.Diff(baseYaml, targetYaml)
		if err != nil {
			diffMap.Results[kDir] = &DiffError{err}
			continue
		}
		diffMap.Results[kDir] = &DiffContent{content}
	}
	return diffMap, nil
}

func MakeBuildOptions(kustomizeLoadRestrictor string) (*krusty.Options, error) {
	var err error
	options := krusty.MakeDefaultOptions()
	if kustomizeLoadRestrictor == "" {
		return options, err
	}
	switch kustomizeLoadRestrictor {
	case "LoadRestrictionsUnknown":
		{
			options.LoadRestrictions = types.LoadRestrictionsUnknown
		}
	case "LoadRestrictionsRootOnly":
		{
			options.LoadRestrictions = types.LoadRestrictionsRootOnly
		}
	case "LoadRestrictionsNone":
		{
			options.LoadRestrictions = types.LoadRestrictionsNone
		}
	default:
		{
			err := errors.Errorf("unknown LoadRestrictions type given by kustomizeLoadRestrictor: %q", kustomizeLoadRestrictor)
			return nil, err
		}
	}
	return options, err
}

type BuildOpts struct {
	KustomizePath           string
	KustomizeLoadRestrictor string
}

func Build(dirPath string, opts BuildOpts) (string, error) {
	if opts.KustomizePath != "" {
		buildArgs := []string {"build"}
		if (opts.KustomizeLoadRestrictor != "") {
			buildArgs = append(buildArgs, "--load-restrictor")
			buildArgs = append(buildArgs, opts.KustomizeLoadRestrictor)
		}
		buildArgs = append(buildArgs, dirPath)
		stdout, _, err := (&utils.WorkDir{}).RunCommand(opts.KustomizePath, buildArgs...)
		if err != nil {
			return "", err
		}
		return stdout, nil
	}
	options, err := MakeBuildOptions(opts.KustomizeLoadRestrictor)
	if err != nil {
		return "", errors.WithStack(err)
	}
	k := krusty.MakeKustomizer(
		options,
	)
	fSys := filesys.MakeFsOnDisk()
	resMap, err := k.Run(fSys, dirPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	bs, err := resMap.AsYaml()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(bs), nil
}
