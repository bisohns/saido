package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	// Verbose : Should display verbose logs
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "saido",
	Short: "TUI for monitoring specific host metrics",
	Long:  `Saido`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Log only errors except in Verbose mode
		if verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n\nSaido - Bisoncorp (2020) (https://github.com/bisoncorps/saido)")
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
	//  cobra.OnInitialize(initConfig)
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}
