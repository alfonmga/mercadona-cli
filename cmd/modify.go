package cmd

import (
	"fmt"

	"github.com/alfonmga/mercadona-cli/internal/pkg/mercadona"
	"github.com/spf13/cobra"
)

var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Get URL to modify active order",
	Long:  `Get URL to modify active order`,
	Run: func(cmd *cobra.Command, args []string) {
		url := mercadona.GetActiveOrderModifyURL()
		fmt.Println(url)
	},
}

func init() {
	orderCmd.AddCommand(modifyCmd)
}
