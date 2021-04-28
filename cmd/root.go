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

/* func init() {
cobra.OnInitialize(initConfig)
rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mercadona-cli.yaml)")
rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
} */

/* func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".mercadona-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
} */
