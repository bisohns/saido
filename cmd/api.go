/*
Copyright © 2021 Bisohns

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
package cmd

import (
	"net/http"
	"os"
	"strconv"

	"github.com/bisohns/saido/client"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	port   string
	server = http.NewServeMux()
)

var apiCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Run saido dashboard on a PORT env variable, fallback to set argument",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		hosts := client.NewHostsController(cfg)
		server.Handle("/metrics", hosts)
		log.Info("listening on :", port)
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal(err)
		}
		go hosts.Run()
		loggedRouters := handlers.LoggingHandler(os.Stdout, server)
		if err := http.ListenAndServe(":"+port, loggedRouters); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	apiCmd.Flags().StringVarP(&port, "port", "p", "3000", "Port to run application server on")
	rootCmd.AddCommand(apiCmd)
}
