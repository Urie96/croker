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
	"time"

	"github.com/gosuri/uitable"
	"github.com/hako/durafmt"
	"github.com/spf13/cobra"
	"github.com/urie96/croker/croker/job"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List jobs",
	Run: func(cmd *cobra.Command, args []string) {
		table := uitable.New()
		table.MaxColWidth = 20

		table.AddRow("ID", "COMMAND", "SCHEDULE", "CREATED", "STATUS")
		// durafmt.ParseShort(time.Second*())
		for _, j := range job.Infos() {
			createdAt := durafmt.ParseShort(time.Since(time.Unix(j.CreatedAt, 0))).String() + " ago"
			var status string
			if j.IsRunning {
				status = "Up " + durafmt.ParseShort(time.Since(time.Unix(j.StartedAt, 0))).String()
			} else {
				status = "Stopped " + durafmt.ParseShort(time.Since(time.Unix(j.StoppedAt, 0))).String()
			}
			table.AddRow(j.ID, j.Command, j.CronSpec, createdAt, status)
		}
		fmt.Println(table)
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(psCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// psCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// psCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
