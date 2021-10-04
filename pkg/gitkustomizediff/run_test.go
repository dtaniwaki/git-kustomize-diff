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
	"strings"
	"testing"

	"github.com/dtaniwaki/git-kustomize-diff/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var (
	exampleRepoUrl = "https://github.com/dtaniwaki/git-kustomize-diff-example.git"
)

func TestRun(t *testing.T) {
	expectedFooDiff := strings.TrimLeft(`
@@ -5,4 +5,4 @@
 spec:
   containers:
   - image: nginx:latest
-    name: foo-modified
+    name: foo-in-branch
`, "\n")
	tmpGitDir, err := ioutil.TempDir("", "kustomize-diff-test-")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer os.RemoveAll(tmpGitDir)
	workDir := &utils.WorkDir{
		Dir: tmpGitDir,
	}
	_, _, err = workDir.RunCommand("git", "clone", exampleRepoUrl, tmpGitDir)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	res, err := Run(tmpGitDir, RunOpts{
		Base:   "origin/main",
		Target: "origin/a-branch",
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, "6206e0c", res.BaseCommit)
	assert.Equal(t, "5a1c160", res.TargetCommit)
	assert.Equal(t, []string{"foo"}, res.DiffMap.Dirs())
	assert.Equal(t, expectedFooDiff, res.DiffMap.Results["foo"].ToString())
}
