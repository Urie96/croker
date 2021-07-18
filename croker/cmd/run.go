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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urie96/croker/croker/job"
)

var detach = true
var schedule string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command or a cronjob in this host",
	Run: func(cmd *cobra.Command, args []string) {
		arg := ""
		for _, a := range args {
			arg += " " + a
		}
		arg = strings.TrimSpace(arg)
		id := job.Add(schedule, arg)
		if detach {
			fmt.Println(id)
		} else {
			job.Log(id, true)
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "crontab expression: * * * * *")
	// runCmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run job in background and print job ID")
}
