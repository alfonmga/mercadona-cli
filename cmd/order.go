package cmd

import (
	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Manage orders",
	Long:  `Manage orders`,
}

func init() {
	rootCmd.AddCommand(orderCmd)
}
