package main

import "github.com/spf13/cobra"

// logout represents the version command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logouts to gnote.",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}
