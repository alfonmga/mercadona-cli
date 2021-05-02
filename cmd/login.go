package cmd

import (
	"github.com/alfonmga/mercadona-cli/internal/pkg/mercadona"
	"github.com/spf13/cobra"
)

func init() {
	var email, password string

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Account login",
		Long:  `Account login`,
		Run: func(cmd *cobra.Command, args []string) {
			mercadona.Authenticate(mercadona.MercadonaLogInCredentialsBodyRequest{Email: email, Password: password})
		},
	}

	loginCmd.Flags().StringVarP(&email, "email", "u", "", "")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")

	authCmd.AddCommand(loginCmd)
}
