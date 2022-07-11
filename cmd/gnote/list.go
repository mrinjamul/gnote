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
	ID string
)

// listCmd represents the version command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"view", "read"},
	Short:   "list note(s).",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			ID = args[0]
		}

		// Read token from config
		config, err := utils.GetConfig()
		if err != nil {
			panic(err)
		}

		if ID != "" {
			note, err := utils.GetNote(ID, config.Token)
			if err != nil {
				fmt.Println(err)
				return
			}
			utils.PrintNote(note)
			return
		}

		// get all notes and print them
		notes, err := utils.GetNotes(config.Token)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, note := range notes {
			utils.PrintNote(note)
		}
	},
}

func init() {
}
