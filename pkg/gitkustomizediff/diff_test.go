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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	wd, _ := os.Getwd()

	expectedYaml := strings.TrimLeft(`
apiVersion: v1
kind: Pod
metadata:
  name: sub1
spec:
  containers:
  - image: nginx:latest
    name: sub1
`, "\n")

	fixturesDirPath := filepath.Join(wd, "fixtures", "diff", "base", "sub1")
	actualYaml, err := Build(fixturesDirPath)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, expectedYaml, actualYaml)
}

func TestDiff(t *testing.T) {
	wd, _ := os.Getwd()

	expectedSub1Diff := strings.TrimLeft(`
@@ -5,4 +5,4 @@
 spec:
   containers:
   - image: nginx:latest
-    name: sub1
+    name: sub1-modified
`, "\n")
	expectedSub2Diff := ""
	expectedInvalidErrorRegexp, _ := regexp.Compile("^accumulating resources: accumulation")

	baseDirPath := filepath.Join(wd, "fixtures", "diff", "base")
	targetDirPath := filepath.Join(wd, "fixtures", "diff", "target")
	diffMap, err := Diff(baseDirPath, targetDirPath, DiffOpts{})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 3, len(diffMap.Results))
	assert.Equal(t, expectedSub1Diff, diffMap.Results["sub1"].(*DiffContent).ToString())
	assert.Equal(t, expectedSub2Diff, diffMap.Results["sub2"].(*DiffContent).ToString())
	assert.Regexp(t, expectedInvalidErrorRegexp, diffMap.Results["invalid"].(*DiffError).Error().Error())
}
