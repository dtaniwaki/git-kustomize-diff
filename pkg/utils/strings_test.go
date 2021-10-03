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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	expectedDiff := strings.TrimLeft(`
@@ -1,3 +1,3 @@
 a
-b
 c
+d
`, "\n")

	diff, err := Diff("a\nb\nc\n", "a\nc\nd\n")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, expectedDiff, diff)

	diff, err = Diff("a\nb\nc\n", "a\nb\nc\n")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, "", diff)
}
