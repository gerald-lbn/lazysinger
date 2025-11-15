package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Environment string

const (
	// EnvDevelopment represents the development environment.
	EnvDevelopment Environment = "development"

	// EnvProduction represents the production environment.
	EnvProduction Environment = "production"

	// EnvTest represents the test environment.
	EnvTest Environment = "test"

	// EnvironmentVariableName is the name of the environment variable used to set the application environment.
	AppEnvironmentVariableName = "REFRAIN_APP_ENVIRONMENT"
)

// SwitchEnvironment sets the environment variable used to dictate which environment the application is
// currently running in.
// This must be called prior to loading the configuration in order for it to take effect.
func SwitchEnvironment(env Environment) {
	if err := os.Setenv(AppEnvironmentVariableName, string(env)); err != nil {
		panic(err)
	}
}

type (
	// Config stores complete application configuration.
	Config struct {
		App       AppConfig
		Database  DatabaseConfig
		HTTP      HTTPConfig
		Libraries LibrariesConfig
		Redis     RedisConfig
		Tasks     TasksConfig
	}

	// AppConfig stores the application configuration.
	AppConfig struct {
		Name          string
		Environment   Environment
		EncryptionKey string
		Timeout       time.Duration
	}

	// DatabaseConfig stores the database configuration.
	DatabaseConfig struct {
		Driver         string
		Connection     string
		TestConnection string
	}

	// HTTPConfig stores HTTP configuration.
	HTTPConfig struct {
		Hostname        string
		Port            uint16
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
		TLS             struct {
			Enabled     bool
			Certificate string
			Key         string
		}
	}

	// LibrariesConfig stores configuration for music libraries.
	LibrariesConfig struct {
		Paths []string `mapstructure:"paths"`
	}

	// RedisConfig stores configuration for redis
	RedisConfig struct {
		Addr string
	}

	// TasksConfig stores the tasks configuration.
	TasksConfig struct {
		GoRoutines      int
		ReleaseAfter    time.Duration
		CleanupInterval time.Duration
		ShutdownTimeout time.Duration
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
