package config

import "github.com/spf13/viper"

var cfg *config

type config struct {
	API            APIConfig
	Authentication AuthConfig
	ModPacks       ModPacksConfig
	User           UserConfig
}

type APIConfig struct {
	Port         string
	ApiKeyHeader string
	Timeout      string
}

type AuthConfig struct {
	Secret     string
	Expiration int64
	AesKey     string
}

type ModPacksConfig struct {
	PublicPath   string
	ManifestName string
}

type UserConfig struct {
	MaxApiKeys int64
}

func init() {
	viper.SetDefault("api.port", "8000")
	viper.SetDefault("api.apikey_header", "x-api-key")

	viper.SetDefault("modpacks.publicpath", "public/modpacks")
	viper.SetDefault("modpacks.manifest_name", "manifest.json")
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
	cfg.API = APIConfig{
		Port:         viper.GetString("api.port"),
		ApiKeyHeader: viper.GetString("api.apikey_header"),
	}
	cfg.ModPacks = ModPacksConfig{
		PublicPath:   viper.GetString("modpacks.publicpath"),
		ManifestName: viper.GetString("modpacks.manifest_name"),
	}

	cfg.Authentication = AuthConfig{
		Secret:     viper.GetString("authentication.secret"),
		Expiration: viper.GetInt64("authentication.expiration"),
		AesKey:     viper.GetString("authentication.aeskey"),
	}
	cfg.User = UserConfig{
		MaxApiKeys: viper.GetInt64("user.maxapikeys"),
	}
	return nil
}

func GetConfig() *config {
	return cfg
}

func GetAesSecret() []byte {
	return []byte(cfg.Authentication.AesKey)
}
func GetJWTSecret() []byte {
	return []byte(cfg.Authentication.Secret)
}
