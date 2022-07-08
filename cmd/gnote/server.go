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
	"embed"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrinjamul/gnote/api/routes"
	"github.com/spf13/cobra"
)

//go:embed views/*
var viewsFs embed.FS

var (
	// startTime is the time when the server starts
	startTime time.Time = time.Now()
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		// server
		// Get port from env
		port := ":8080"
		_, present := os.LookupEnv("PORT")
		if present {
			port = ":" + os.Getenv("PORT")

		}
		// Set the router as the default one shipped with Gin
		server := gin.Default()
		// Initialize the routes
		routes.StartTime = startTime
		routes.ViewsFs = viewsFs
		routes.InitRoutes(server)
		routes.BootTime = time.Since(startTime)
		// Start and run the server
		log.Fatal(server.Run(port))
	},
}
