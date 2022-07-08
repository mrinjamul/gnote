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
	flagID int
)

// listCmd represents the version command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list note(s).",
	Run: func(cmd *cobra.Command, args []string) {
		if flagID != -1 {
			data, err := utils.GetNote(fmt.Sprintf("%d", flagID))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(data)
			return
		}
		// get all notes and print them
		notes, err := utils.GetNotes()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(notes)
	},
}

func init() {
	listCmd.Flags().IntVarP(&flagID, "view", "v", -1, "id of the note")
}
