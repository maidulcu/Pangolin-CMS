package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	SiteURL  string `mapstructure:"site_url"`
	APIKey   string `mapstructure:"api_key"`
	S3Bucket string `mapstructure:"s3_bucket"`
	S3Region string `mapstructure:"s3_region"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("pangolin")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.pangolin")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found, run 'pangolin init' first")
		}
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	viper.Set("site_url", cfg.SiteURL)
	viper.Set("api_key", cfg.APIKey)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := homeDir + "/.pangolin"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	viper.SetConfigName("pangolin")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.SafeWriteConfig(); err != nil {
		return viper.WriteConfig()
	}

	return nil
}
