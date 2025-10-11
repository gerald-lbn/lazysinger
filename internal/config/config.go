package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MusicLibraryPath string
	RedisAddr        string
}

const (
	// Environment variable names
	MUSIC_LIBRARY_PATH = "MUSIC_LIBRARY_PATH"
	REDIS_ADDR         = "REDIS_ADDR"

	// Default values
	DEFAULT_MUSIC_LIBRARY_PATH = "/music"
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

	return &Config{
		MusicLibraryPath: musicPath,
		RedisAddr:        redisAddr,
	}
}
