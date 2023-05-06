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
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
)

func TestMakeBuildOptions(t *testing.T) {
	var err error
	var kustomizeLoadRestrictor string
	var options *krusty.Options
	defaultOptions := krusty.MakeDefaultOptions()

	kustomizeLoadRestrictor = ""
	options, err = MakeBuildOptions(kustomizeLoadRestrictor)
	assert.Equal(t, err, nil)
	assert.Equal(t, options, defaultOptions)

	kustomizeLoadRestrictor = "LoadRestrictionsUnknown"
	options, err = MakeBuildOptions(kustomizeLoadRestrictor)
	assert.Equal(t, err, nil)
	assert.Equal(t, options.LoadRestrictions, types.LoadRestrictionsUnknown)

	kustomizeLoadRestrictor = "LoadRestrictionsRootOnly"
	options, err = MakeBuildOptions(kustomizeLoadRestrictor)
	assert.Equal(t, err, nil)
	assert.Equal(t, options.LoadRestrictions, types.LoadRestrictionsRootOnly)

	kustomizeLoadRestrictor = "LoadRestrictionsNone"
	options, err = MakeBuildOptions(kustomizeLoadRestrictor)
	assert.Equal(t, err, nil)
	assert.Equal(t, options.LoadRestrictions, types.LoadRestrictionsNone)

	invalidType := "invalid-load-restrictions-type--"
	kustomizeLoadRestrictor = invalidType
	options, err = MakeBuildOptions(kustomizeLoadRestrictor)
	assert.Equal(t, options, (*krusty.Options)(nil))
	assert.NotEqual(t, err, nil)
	assert.Error(t, err, "unknown LoadRestrictions type given by kustomizeLoadRestrictor: %q", invalidType)
}

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
	actualYaml, err := Build(fixturesDirPath, BuildOpts{})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, expectedYaml, actualYaml)
}

func TestBuildLoadRestrictionsNone(t *testing.T) {
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

	fixturesDirPath := filepath.Join(wd, "fixtures", "diff-load-restrictions-none", "base", "sub1/nested")
	_, err := Build(fixturesDirPath, BuildOpts{})
	assert.NotEqual(t, err, nil)

	buildOpts := BuildOpts{"", "LoadRestrictionsNone"}
	actualYaml, err := Build(fixturesDirPath, buildOpts)
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

func TestDiffLoadRestrictionsNone(t *testing.T) {
	wd, _ := os.Getwd()

	expectedSub1Diff := strings.TrimLeft(`
@@ -5,4 +5,4 @@
 spec:
   containers:
   - image: nginx:latest
-    name: sub1
+    name: sub1-modified
`, "\n")

	baseDirPath := filepath.Join(wd, "fixtures", "diff-load-restrictions-none", "base")
	targetDirPath := filepath.Join(wd, "fixtures", "diff-load-restrictions-none", "target")
	_, err := Diff(baseDirPath, targetDirPath, DiffOpts{})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	diffOpts := DiffOpts{nil, nil, "", "LoadRestrictionsNone"}
	diffMap, err := Diff(baseDirPath, targetDirPath, diffOpts)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, 1, len(diffMap.Results))
	assert.Equal(t, expectedSub1Diff, diffMap.Results["sub1/nested"].(*DiffContent).ToString())
}
