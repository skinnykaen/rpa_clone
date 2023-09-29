package configs

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/spf13/viper"
)

func Init(m consts.Mode) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	switch m {
	case consts.Production:
		viper.SetConfigName("production")
		viper.SetConfigType("env")
		viper.AddConfigPath(viper.GetString("production_env.path"))
	case consts.Development:
		viper.SetConfigName("development")
		viper.SetConfigType("env")
		viper.AddConfigPath(viper.GetString("development_env.path"))
	}
	return viper.MergeInConfig()
}
