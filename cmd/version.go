package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get saido version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
