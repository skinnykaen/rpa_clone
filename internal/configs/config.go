package configs

import (
	"github.com/spf13/viper"
	"rpa_clone/internal/consts"
)

func Init(m consts.Mode) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../internal/configs/")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	switch m {
	case consts.Production:
		viper.SetConfigName("production")
		viper.SetConfigType("env")
		viper.AddConfigPath(viper.GetString("production_env.path"))
		break
	case consts.Development:
		viper.SetConfigName("development")
		viper.SetConfigType("env")
		viper.AddConfigPath(viper.GetString("development_env.path"))
		break
	}
	err = viper.MergeInConfig()
	if err != nil {
		return err
	}
	return nil
}
