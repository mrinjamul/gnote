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

// removeCmd represents the version command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"del", "rm", "delete"},
	Short:   "remove a note.",
	Run: func(cmd *cobra.Command, args []string) {
		// parse args
		if len(args) == 0 {
			fmt.Println("error: too short argument")
			fmt.Println("Usage: gnote remove [id]")
			return
		}
		id := args[0]

		// Read token from config
		config, err := utils.GetConfig()
		if err != nil {
			panic(err)
		}

		note, err := utils.DeleteNote(id, config.Token)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d: %s has been deleted!\n", note.ID, note.Title)
	},
}
