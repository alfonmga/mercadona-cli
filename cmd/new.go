package cmd

import (
	"github.com/alfonmga/mercadona-cli/internal/pkg/mercadona"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Make new order",
	Long:  `Make new order`,
	Run: func(cmd *cobra.Command, args []string) {
		mercadona.MakeNewOrder()
	},
}

func init() {
	orderCmd.AddCommand(newCmd)
}
