package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort              string `mapstructure:"APP_PORT"`
	AwsRegion            string `mapstructure:"AWS_REGION"`
	LeapcellBaseEndpoint string `mapstructure:"LEAPCELL_BASE_ENDPOINT"`
	LeapcellCdn          string `mapstructure:"LEAPCELL_CDN"`
	AwsAccessKeyId       string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey   string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsBucket            string `mapstructure:"AWS_BUCKET"`
}

func LoadConfig() (Config, error) {
	viper.AutomaticEnv()

	var conf Config

	err := bindEnv(conf)
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return Config{}, errors.New("failed to parse configuration into struct")
	}

	return conf, nil
}

func bindEnv(cnf Config) error {

	var errs []error

	r := reflect.TypeOf(cnf)

	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		tag := field.Tag.Get("mapstructure")

		if tag == "" {
			continue
		}

		_, exists := os.LookupEnv(tag)
		if !exists {
			errs = append(errs, fmt.Errorf("%v is missing in environment variable", tag))
		}

		err := viper.BindEnv(tag)
		if err != nil {
			errs = append(errs, err)

		}

	}

	return errors.Join(errs...)
}
