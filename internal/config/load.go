package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Load(env string) (Config, error) {
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/") // path to look for the config file

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	// enable env vars for secrets and overrides
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	cfg.Env = env

	return cfg, nil
}
