package cmd

import (
	"github.com/alfonmga/mercadona-cli/internal/pkg/mercadona"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Current account information",
	Long:  `Current account information`,
	Run: func(cmd *cobra.Command, args []string) {
		mercadona.CustomerInfo()
	},
}

func init() {
	authCmd.AddCommand(infoCmd)
}
