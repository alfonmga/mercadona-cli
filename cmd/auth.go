package cmd

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication",
	Long:  `Authentication`,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
