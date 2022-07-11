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
	"strings"

	"github.com/mrinjamul/gnote/models"
	"github.com/mrinjamul/gnote/utils"
	"github.com/spf13/cobra"
)

// searchCmd represents the version command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search note(s).",
	Run: func(cmd *cobra.Command, args []string) {
		// parse args
		if len(args) == 0 {
			fmt.Println("error: too short argument")
			fmt.Println("Usage: gnote search [content]")
			return
		}
		content := args[0]
		// Read token from config
		config, err := utils.GetConfig()
		if err != nil {
			panic(err)
		}
		// get all notes and print them
		notes, err := utils.GetNotes(config.Token)
		if err != nil {
			fmt.Println(err)
			return
		}
		var foundNotes []models.Note
		for _, note := range notes {
			if strings.Contains(note.Title, content) || strings.Contains(note.Content, content) {
				foundNotes = append(foundNotes, note)
			}
		}

		// Print notes
		fmt.Printf("%d notes found: \n", len(foundNotes))
		for _, note := range foundNotes {
			utils.PrintNote(note)
		}
	},
}
