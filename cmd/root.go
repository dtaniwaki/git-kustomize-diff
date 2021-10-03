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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	verbose int
}

var RootCmd = &cobra.Command{
	Use:           "git-kustomize-diff",
	Short:         "git-kustomize-diff",
	SilenceErrors: false,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use == "version" {
			return
		}
		switch rootOpts.verbose {
		case 0:
			log.SetLevel(log.ErrorLevel)
		case 1:
			log.SetLevel(log.InfoLevel)
		case 2:
			log.SetLevel(log.DebugLevel)
		case 3:
			log.SetLevel(log.TraceLevel)
		default:
			log.SetLevel(log.TraceLevel)
		}
	},
}

var rootOpts rootFlags

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().CountVarP(&rootOpts.verbose, "verbose", "v", "verbose mode. (1: info, 2: debug, 3: trace)")
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(runCmd)
}
