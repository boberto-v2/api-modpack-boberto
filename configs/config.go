package config

import "github.com/spf13/viper"

type config struct {
	Port       string
	MaxApiKeys int64
	PublicPath string
}

var cfg *config

func init() {
	viper.SetDefault("api.port", 8000)
	viper.SetDefault("api.publicpath", "public/modpacks")
}

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	cfg = new(config)
	cfg.Port = viper.GetString("api.port")
	cfg.PublicPath = viper.GetString("api.publicpath")
	return nil
}

func GetConfig() *config {
	return cfg
}
