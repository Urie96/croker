/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/urie96/croker/crokerd/consts"
	"github.com/urie96/croker/utils"
)

var follow bool

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Fetch the logs of the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		var r io.ReadCloser
		var err error
		if follow {
			r, err = utils.FollowFile(consts.DaemonLogPath)
		} else {
			r, err = os.Open(consts.DaemonLogPath)
		}
		logFatal(err)
		io.Copy(os.Stdout, r)
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "")
}
