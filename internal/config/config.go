package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MusicLibraryPath     string
	RedisAddr            string
	WorkerConcurrency    int
	WorkerStrictPriority bool
}

const (
	// Environment variable names
	MUSIC_LIBRARY_PATH     = "MUSIC_LIBRARY_PATH"
	REDIS_ADDR             = "REDIS_ADDR"
	WORKER_CONCURRENCY     = "WORKER_CONCURRENCY"
	WORKER_STRICT_PRIORITY = "WORKER_STRICT_PRIORITY"

	// Default values
	DEFAULT_MUSIC_LIBRARY_PATH     = "/music"
	DEFAULT_WORKER_CONCURRENCY     = "10"
	DEFAULT_WORKER_STRICT_PRIORITY = "true"
)

func GetEnv(name string) string {
	variable := os.Getenv(name)
	if variable == "" {
		log.Fatalf("%s environment variable not set", name)
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

	musicPath := GetEnvWithDefault(MUSIC_LIBRARY_PATH, DEFAULT_MUSIC_LIBRARY_PATH)
	redisAddr := GetEnv(REDIS_ADDR)
	workerConcurrency := GetEnvWithDefault(WORKER_CONCURRENCY, DEFAULT_WORKER_CONCURRENCY)
	workerConcurrencyInt, err := strconv.Atoi(workerConcurrency)
	if err != nil {
		log.Fatalf("Invalid value for %s: %v", WORKER_CONCURRENCY, err)
	}
	workerStrictPriority := GetEnvWithDefault(WORKER_STRICT_PRIORITY, DEFAULT_WORKER_STRICT_PRIORITY)
	workerStrictPriorityBool, err := strconv.ParseBool(workerStrictPriority)
	if err != nil {
		log.Fatalf("Invalid value for %s: %v", WORKER_STRICT_PRIORITY, err)
	}

	return &Config{
		MusicLibraryPath:     musicPath,
		RedisAddr:            redisAddr,
		WorkerConcurrency:    workerConcurrencyInt,
		WorkerStrictPriority: workerStrictPriorityBool,
	}
}
