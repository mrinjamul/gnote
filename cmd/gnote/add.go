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

	"github.com/mrinjamul/gnote/utils"
	"github.com/spf13/cobra"
)

var (
	flagTitle   string
	flagContent string
)

// addCmd represents the version command
var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Short:   "add a note.",
	Run: func(cmd *cobra.Command, args []string) {
		if flagTitle == "" || flagContent == "" {
			fmt.Println("title and content are required")
			return
		}

		// Read token from config
		config, err := utils.GetConfig()
		if err != nil {
			panic(err)
		}

		fmt.Println("Title:", flagTitle)
		fmt.Println("Body:", flagContent)
		data, err := utils.CreateNote(flagTitle, flagContent, config.Token)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
	},
}

func init() {
	addCmd.Flags().StringVarP(&flagTitle, "title", "t", "", "title of the note")
	addCmd.Flags().StringVarP(&flagContent, "content", "c", "", "content of the note")
}
