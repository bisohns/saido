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
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bisohns/saido/client"
	"github.com/bisohns/saido/config"
	"github.com/gorilla/handlers"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const appName = "saido"

var (
	port   string
	server = http.NewServeMux()
	// Verbose : Should display verbose logs
	verbose bool
	// open browser
	browserFlag bool

	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "saido",
	Short: "Tool for monitoring metrics",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Log only errors except in Verbose mode
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n\nSaido - Bisohns (2020) (https://github.com/bisohns/saido)")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		fmt.Println(args)
		cfg = config.LoadConfig(cfgFile)
		hosts := client.NewHostsController(cfg)

		server.Handle("/metrics", hosts)
		log.Info("listening on :", port)
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal(err)
		}
		go hosts.Run()
		loggedRouters := handlers.LoggingHandler(os.Stdout, server)
		// Trigger browser open
		url := fmt.Sprintf("http://localhost:%s", port)
		if browserFlag {
			browser.OpenURL(url)
		}
		if err := http.ListenAndServe(":"+port, loggedRouters); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&port, "port", "p", "3000", "Port to run application server on")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Run saido in verbose mode")
	rootCmd.Flags().BoolVarP(&browserFlag, "open-browser", "b", false, "Prompt open browser")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to config file")
	if len(os.Args) >= 2 && os.Args[1] != "version" || len(os.Args) == 1 {
		cobra.MarkFlagRequired(rootCmd.PersistentFlags(), "config")
	}
}
