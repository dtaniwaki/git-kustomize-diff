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

	"github.com/dtaniwaki/git-kustomize-diff/pkg/utils"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func Diff(baseDirPath, targetDirPath string) (*DiffMap, error) {
	log.Info("Start diff")
	baseKDirs, err := utils.ListKustomizeDirs(baseDirPath, utils.ListKustomizeDirsOpts{})
	if err != nil {
		return nil, err
	}
	targetKDirs, err := utils.ListKustomizeDirs(targetDirPath, utils.ListKustomizeDirsOpts{})
	if err != nil {
		return nil, err
	}
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
		baseYaml, err := Build(baseKDirPath)
		if err != nil {
			diffMap.Results[kDir] = &DiffError{err}
			continue
		}
		targetYaml, err := Build(targetKDirPath)
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

func Build(dirPath string) (string, error) {
	fSys := filesys.MakeFsOnDisk()
	k := krusty.MakeKustomizer(
		krusty.MakeDefaultOptions(),
	)
	resMap, err := k.Run(fSys, dirPath)
	if err != nil {
		return "", err
	}
	bs, err := resMap.AsYaml()
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
