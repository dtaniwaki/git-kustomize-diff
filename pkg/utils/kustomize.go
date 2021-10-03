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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type ListKustomizeDirsOpts struct {
	IncludeRegexp *regexp.Regexp
	ExcludeRegexp *regexp.Regexp
}

func ListKustomizeDirs(dirPath string, opts ListKustomizeDirsOpts) ([]string, error) {
	targetFiles := make([]string, 0)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if !KustomizationExists(path) {
			return nil
		}
		included := true
		if opts.IncludeRegexp != nil {
			m := opts.IncludeRegexp.Match([]byte(path))
			if !m {
				included = false
			}
		}
		if included {
			if opts.ExcludeRegexp != nil {
				m := opts.ExcludeRegexp.Match([]byte(path))
				if m {
					included = false
				}
			}
		}
		if included {
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			targetFiles = append(targetFiles, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return targetFiles, nil
}

func KustomizationExists(path string) bool {
	return Exists(filepath.Join(path, "kustomization.yaml")) || Exists(filepath.Join(path, "kustomization.yml"))
}

func MakeKustomizeDir(dirPath string) error {
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		return err
	}
	kustomizationFilePath := filepath.Join(dirPath, "kustomization.yaml")
	if Exists(kustomizationFilePath) {
		return fmt.Errorf("File already exists: %s", kustomizationFilePath)
	}
	f, err := os.Create(kustomizationFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
