package cmd

import (
	"github.com/spf13/cobra"
)

// var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "mercadona-cli",
	Short: "Mercadona CLI",
	Long:  `Mercadona CLI`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
