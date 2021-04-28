package cmd

import (
	"github.com/alfonmga/mercadona-cli/internal/pkg/mercadona"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Account login",
	Long:  `Account login`,
	Run: func(cmd *cobra.Command, args []string) {
		mercadona.Authenticate(mercadona.MercadonaLogInCredentialsBodyRequest{Email: "<email>", Password: "<password>"})
	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
