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
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bisohns/saido/config"
)

var (
	// Verbose : Should display verbose logs
	verbose       bool
	dashboardInfo *config.DashboardInfo
)

const appName = "saido"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "saido",
	Short: "Tool for monitoring metrics",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Log only errors except in Verbose mode
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n\nSaido - Bisohns (2020) (https://github.com/bisohns/saido)")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Run dashboard with the [dashboard] subcommand. Feel free to drop a star at https://github.com/bisohns/saido")
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
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Run saido in verbose mode")
}
