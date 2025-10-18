package config

import (
	"os"
	"strconv"

	"github.com/gerald-lbn/lazysinger/log"
)

type configOptions struct {
	Logger  loggerOptions
	Scanner scannerOptions
	Worker  workerOptions
}

type loggerOptions struct {
	Level string
}

type scannerOptions struct {
	Directory string
	Schedule  string
}

type workerOptions struct {
	Concurrency int
	RedisAddr   string
}

var (
	// Environement keys
	LOG_LEVEL          = "LOG_LEVEL"
	MUSIC_LIBRARY      = "MUSIC_LIBRARY"
	REDIS_ADDR         = "REDIS_ADDR"
	SCHEDULE           = "SCHEDULE"
	WORKER_CONCURRENCY = "WORKER_CONCURRENCY"

	// Default environment keys
	DEFAULT_LOG_LEVEL          = "info"
	DEFAULT_MUSIC_LIBRARY      = "/music"
	DEFAULT_REDIS_ADDR         = "localhost:6379"
	DEFAULT_SCHEDULE           = "* */1 * * *"
	DEFAULT_WORKER_CONCURRENCY = "1"

	Server = &configOptions{}
)

// LoadWithDefault checks for an environment variable and returns its value
// or a default value if not set.
func LoadWithDefault(key string, defaultValue string) string {
	val, present := os.LookupEnv(key)
	log.Debug().Str("key", key).Str("val", val).Str("default", defaultValue).Send()
	if !present {
		return defaultValue
	}
	return val
}

// Setup loads the configuration from environment variables,
// applying default values where not specified.
// This function is now exported to be called explicitly by tests.
func Setup() error {
	Server.Logger.Level = LoadWithDefault(LOG_LEVEL, DEFAULT_LOG_LEVEL)

	Server.Scanner.Directory = LoadWithDefault(MUSIC_LIBRARY, DEFAULT_MUSIC_LIBRARY)
	Server.Scanner.Schedule = LoadWithDefault(SCHEDULE, DEFAULT_SCHEDULE)

	c := LoadWithDefault(WORKER_CONCURRENCY, DEFAULT_WORKER_CONCURRENCY)
	i, err := strconv.Atoi(c)
	if err != nil || i <= 0 {
		log.Fatal().Err(err).Str("key", WORKER_CONCURRENCY).Str("value", c).Msg("Invalid value for environement key")
		return err
	}

	Server.Worker.Concurrency = i
	Server.Worker.RedisAddr = LoadWithDefault(REDIS_ADDR, DEFAULT_REDIS_ADDR)

	return nil
}

// ResetConfig resets the global configuration.
func ResetConfig() {
	Server = &configOptions{} // Reinitialize the Server pointer
}
