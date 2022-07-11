package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/mrinjamul/gnote/utils"
	"github.com/spf13/cobra"
)

// signup represents the version command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "signup to gnote.",
	Run: func(cmd *cobra.Command, args []string) {
		// prompt for username and password
		prompt := promptui.Prompt{
			Label: "Username",
			Validate: func(input string) error {
				ok := utils.IsValidUserName(input)
				if !ok {
					return errors.New("invalid username")
				}
				return nil
			},
		}
		username, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		prompt = promptui.Prompt{
			Label: "Password",
			Mask:  '*',
			Validate: func(input string) error {
				ok := utils.IsValidPassword(input)
				if !ok {
					return errors.New("invalid password")
				}
				return nil
			},
		}
		password, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		// signup
		data, err := utils.CLISignup(username, password)
		if err != nil {
			fmt.Println(err)
			return
		}
		var resp map[string]string
		err = json.Unmarshal(data, &resp)
		if err != nil {
			fmt.Println(err)
			return
		}
		if resp["status"] == "success" {
			fmt.Println("Signup successful")
		} else {
			fmt.Println("Signup failed")
		}
	},
}
