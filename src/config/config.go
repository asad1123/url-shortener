package config

import "github.com/spf13/viper"

type AppConfig struct {
	ServerPort        string `mapstructure:"SERVER_PORT"`
	DbHostname        string `mapstructure:"DB_HOSTNAME"`
	DbPort            string `mapstructure:"DB_PORT"`
	DbName            string `mapstructure:"DB_NAME"`
	CacheHostname     string `mapstructure:"CACHE_HOSTNAME"`
	CachePort         string `mapstructure:"CACHE_PORT"`
	CachePassword     string `mapstructure:"CACHE_PASSWORD"`
	CacheDb           int    `mapstructure:"CACHE_DB"`
	ShortenedIdLength int    `mapstructure:"SHORTENED_ID_LENGTH"`
}

func LoadConfig(path string, configName string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
