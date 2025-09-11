package config

import "github.com/spf13/viper"

type Config struct {
	AppUrl             string `mapstructure:"APP_URL"`
	AwsRegion          string `mapstructure:"AWS_REGION"`
	AwsEndpoint        string `mapstructure:"AWS_ENDPOINT"`
	AwsAccessKeyId     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsBucket          string `mapstructure:"AWS_BUCKET"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var conf Config

	err = viper.Unmarshal(&conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
