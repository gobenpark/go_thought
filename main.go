package main

import (
	"fmt"
	"os"

	"github.com/gobenpark/go_thought/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "gothought",
		Short: "A Proxy for Large Language Model API",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cmd.RunCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	rootCmd.Execute()
}
