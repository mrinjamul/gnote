package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

var EnvConfigs map[string]string
var err error

func LoadEnv() {
	EnvConfigs, err = godotenv.Read()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
}
