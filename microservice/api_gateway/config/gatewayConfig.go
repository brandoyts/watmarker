package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Service struct {
	Name   string `mapstructure:"name"`
	Prefix string `mapstructure:"prefix"`
	Url    string `mapstructure:"url"`
}

type GatewayConfig struct {
	Address  string    `mapstructure:"address"`
	LogLevel string    `mapstructure:"logLevel"`
	Services []Service `mapstructure:"services"`
}

func LoadGatewayConfig(path string) GatewayConfig {
	viper.AddConfigPath(path)
	viper.SetConfigName("gatewayConfig")
	viper.SetConfigType("yml")

	viper.SetConfigName("gatewayConfig")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading gateway configuration: %v", err))
	}

	var cnf GatewayConfig

	err = viper.Unmarshal(&cnf)
	if err != nil {
		panic(fmt.Errorf("error parsing gateway configuration into struct: %v", err))
	}

	return cnf
}
