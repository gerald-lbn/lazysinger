package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Environment string

const (
	// EnvLocal represents the local environment.
	EnvLocal Environment = "local"

	// EnvTest represents the test environment.
	EnvTest Environment = "test"

	// EnvDevelopment represents the development environment.
	EnvDevelopment Environment = "dev"

	// EnvStaging represents the staging environment.
	EnvStaging Environment = "staging"

	// EnvQA represents the qa environment.
	EnvQA Environment = "qa"

	// EnvProduction represents the production environment.
	EnvProduction Environment = "prod"
)

// SwitchEnvironment sets the environment variable used to dictate which environment the application is
// currently running in.
// This must be called prior to loading the configuration in order for it to take effect.
func SwitchEnvironment(env Environment) {
	if err := os.Setenv("REFRAIN_APP_ENVIRONMENT", string(env)); err != nil {
		panic(err)
	}
}

type (
	// Config stores complete application configuration.
	Config struct {
		App AppConfig
	}

	// AppConfig stores the application configuration.
	AppConfig struct {
		Name        string
		Environment Environment
	}
)

// GetConfig loads and returns configuration.
func GetConfig() (Config, error) {
	var cfg Config

	// Load the config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("../../config/")

	// Load environment variables
	viper.SetEnvPrefix("REFRAIN")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
