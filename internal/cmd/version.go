/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "CLI version",
	Long: `Displays version of currently installed gq CLI.

If "dev" is displayed, it means a local build is being used instead of installed release.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("Version: %s\n", version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
