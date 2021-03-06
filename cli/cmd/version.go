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

	"github.com/spf13/cobra"
	"github.com/urie96/croker/cli/remote"
	"github.com/urie96/croker/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the Croker version information",
	Run: func(cmd *cobra.Command, args []string) {
		server := remote.Server{}
		clientV := version.Version
		fmt.Println("Client:", clientV)
		serverV := server.Version()
		fmt.Println("Server:", serverV)
		var cb, cm, sb, sm int32
		fmt.Sscanf(clientV, "%d.%d", &cb, &cm)
		fmt.Sscanf(serverV, "%d.%d", &sb, &sm)
		if cb*100+cm != sb*100+sm {
			fmt.Println("Version mismatch!")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
