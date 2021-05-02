package cmd

import (
	"fmt"

	"github.com/alfonmga/mercadona-cli/internal/pkg/mercadona"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Current account information",
	Long:  `Current account information`,
	Run: func(cmd *cobra.Command, args []string) {
		res := mercadona.CustomerInfo()
		fmt.Println(res)
	},
}

func init() {
	authCmd.AddCommand(infoCmd)
}
