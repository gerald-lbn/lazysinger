package config

import (
	"os"
	"strconv"

	"github.com/gerald-lbn/lazysinger/log"

	"github.com/joho/godotenv"
)

type Config struct {
	MusicLibraryPath     string
	RedisAddr            string
	WorkerConcurrency    int
	WorkerStrictPriority bool
	LogLevel             string
}

const (
	// Environment variable names
	LOG_LEVEL              = "LOG_LEVEL"
	MUSIC_LIBRARY_PATH     = "MUSIC_LIBRARY_PATH"
	REDIS_ADDR             = "REDIS_ADDR"
	WORKER_CONCURRENCY     = "WORKER_CONCURRENCY"
	WORKER_STRICT_PRIORITY = "WORKER_STRICT_PRIORITY"

	// Default values
	DEFAULT_LOG_LEVEL              = "info"
	DEFAULT_MUSIC_LIBRARY_PATH     = "/music"
	DEFAULT_WORKER_CONCURRENCY     = "10"
	DEFAULT_WORKER_STRICT_PRIORITY = "true"
)

func GetEnv(name string) string {
	variable := os.Getenv(name)
	if variable == "" {
		log.Fatal().Str("name", name).Msg("Variable is not set")
	}
	return variable
}

func GetEnvWithDefault(name, defaultValue string) string {
	variable := os.Getenv(name)
	if variable == "" {
		return defaultValue
	}
	return variable
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	logLevel := GetEnvWithDefault(LOG_LEVEL, DEFAULT_LOG_LEVEL)
	musicPath := GetEnvWithDefault(MUSIC_LIBRARY_PATH, DEFAULT_MUSIC_LIBRARY_PATH)
	redisAddr := GetEnv(REDIS_ADDR)
	workerConcurrency := GetEnvWithDefault(WORKER_CONCURRENCY, DEFAULT_WORKER_CONCURRENCY)
	workerConcurrencyInt, err := strconv.Atoi(workerConcurrency)
	if err != nil {
		log.Fatal().Str("key", WORKER_CONCURRENCY).Err(err).Msg("Invalid value")
	}
	workerStrictPriority := GetEnvWithDefault(WORKER_STRICT_PRIORITY, DEFAULT_WORKER_STRICT_PRIORITY)
	workerStrictPriorityBool, err := strconv.ParseBool(workerStrictPriority)
	if err != nil {
		log.Fatal().Str("key", WORKER_STRICT_PRIORITY).Err(err).Msg("Invalid value")
	}

	return &Config{
		LogLevel:             logLevel,
		MusicLibraryPath:     musicPath,
		RedisAddr:            redisAddr,
		WorkerConcurrency:    workerConcurrencyInt,
		WorkerStrictPriority: workerStrictPriorityBool,
	}
}
