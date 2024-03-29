package cmd

import (
	"github.com/spf13/cobra"
)

var file string

var rootCmd = &cobra.Command{
	Use:   "intermediate",
	Short: "An Intermediate application",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&file, "config", "c", "config.yaml", "Config file")

	// sub commands are added in respective files
	return rootCmd.Execute()
}
