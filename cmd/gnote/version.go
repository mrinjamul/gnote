/*
Copyright © 2022 Injamul Mohammad Mollah

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
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Configured via -ldflags during build
// Version is the version of the binary
var Version = "dev"

// GitCommit is the git commit hash of the binary
var GitCommit string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints version.",
	Run: func(cmd *cobra.Command, args []string) {
		shortCommit := shortGitCommit(GitCommit)
		version := fmt.Sprintf("Version: %s %s", Version, shortCommit)
		fmt.Println(version)
	},
}

// shortGitCommit returns the short form of the git commit hash
func shortGitCommit(fullGitCommit string) string {
	shortCommit := ""
	if len(fullGitCommit) >= 6 {
		shortCommit = fullGitCommit[0:6]
	}

	return shortCommit
}
