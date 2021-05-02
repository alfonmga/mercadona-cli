package cmd

import (
	"github.com/alfonmga/mercadona-cli/internal/pkg/mercadona"

	"github.com/spf13/cobra"
)

func init() {
	var date string

	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "Make new order",
		Long:  `Make new order`,
		Run: func(cmd *cobra.Command, args []string) {
			mercadona.MakeNewOrder(date)
		},
	}

	newCmd.Flags().StringVarP(&date, "date", "d", "", "")
	newCmd.MarkFlagRequired("date")

	orderCmd.AddCommand(newCmd)
}
