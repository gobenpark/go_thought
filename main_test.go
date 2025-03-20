package main

import (
	"fmt"
	"testing"

	"github.com/gobenpark/go_thought/internal/proxy"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestViper(t *testing.T) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("gothought")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	p := proxy.Config{}
	err := viper.UnmarshalKey("config", &p)
	require.NoError(t, err)
	fmt.Println(p)
}
