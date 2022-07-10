package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/mrinjamul/gnote/utils"
	"github.com/spf13/cobra"
)

// login represents the version command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to gnote.",
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
		// login
		data, err := utils.CLILogin(username, password)
		if err != nil {
			fmt.Println(err)
			return
		}
		// marshal to models.Config
		var config map[string]string
		err = json.Unmarshal([]byte(data), &config)
		if err != nil {
			fmt.Println(err)
			return
		}
		// save to config file
		err = utils.SaveToken(config["token"])
		if err != nil {
			fmt.Println(err)
			return
		}
		if config["token"] != "" {
			fmt.Println("Login successful")
		} else {
			fmt.Println("Login failed")
		}
	},
}
