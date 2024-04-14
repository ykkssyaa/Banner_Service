package config

import "github.com/spf13/viper"

func InitConfig() error {

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	return err
}
