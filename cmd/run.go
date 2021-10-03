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

package cmd

import (
	"fmt"

	"github.com/dtaniwaki/git-kustomize-diff/pkg/gitkustomizediff"
	"github.com/spf13/cobra"
)

type runFlags struct {
	base   string
	target string
}

var runCmd = &cobra.Command{
	Use:   "run target_dir",
	Short: "Run git-kustomize-diff",
	Long:  `Run git-kustomize-diff`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		opts := gitkustomizediff.RunOpts{
			Base:   runOpts.base,
			Target: runOpts.target,
		}
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}
		err := gitkustomizediff.Run(dir, opts)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var runOpts runFlags

func init() {
	runCmd.PersistentFlags().StringVar(&runOpts.base, "base", "", "base commitish (default to origin/main)")
	runCmd.PersistentFlags().StringVar(&runOpts.target, "target", "", "target commitish (default to the current branch)")
}
