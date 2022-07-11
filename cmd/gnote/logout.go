package main

import (
	"fmt"

	"github.com/mrinjamul/gnote/utils"
	"github.com/spf13/cobra"
)

// logout represents the version command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logouts to gnote.",
	Run: func(cmd *cobra.Command, args []string) {
		// Remove token from config.json
		utils.SaveToken("")
		fmt.Println("Logout sucessfully.")
	},
}
