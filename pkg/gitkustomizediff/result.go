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
	"sort"
	"strings"
)

type DiffResult interface {
	ToString() string
	AsMarkdown() string
}

type DiffError struct {
	err error
}

func (r *DiffError) ToString() string {
	return fmt.Sprintf("%s", r.Error())
}

func (r *DiffError) AsMarkdown() string {
	return fmt.Sprintf("```\n%s```", r.Error())
}

func (r *DiffError) Error() error {
	return r.err
}

type DiffContent struct {
	content string
}

func (r *DiffContent) ToString() string {
	return r.content
}

func (r *DiffContent) AsMarkdown() string {
	if r.content == "" {
		return "```N/A```"
	} else {
		return fmt.Sprintf("```diff\n%s```", r.content)
	}
}

type DiffMap struct {
	Results map[string]DiffResult
}

func NewDiffMap() *DiffMap {
	return &DiffMap{
		Results: make(map[string]DiffResult),
	}
}

func (dm *DiffMap) AsMarkdown() string {
	paths := make([]string, 0)
	for path := range dm.Results {
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool {
		return paths[i] < paths[j]
	})
	lines := make([]string, len(paths))
	for idx, path := range paths {
		lines[idx] = fmt.Sprintf("## %s:\n%s", path, dm.Results[path].AsMarkdown())
	}
	return strings.Join(lines, "\n")
}
