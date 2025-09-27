package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MusicLibraryPath string
	RedisAddr        string
	RedisPassword    string
}

const (
	MUSIC_LIBRARY_PATH = "MUSIC_LIBRARY_PATH"
	REDIS_ADDR         = "REDIS_ADDR"
)

func getEnv(name string) string {
	variable := os.Getenv(name)
	if variable == "" {
		log.Fatalf("%s environment variable not set", name)
	}
	return variable
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	musicPath := getEnv(MUSIC_LIBRARY_PATH)
	redisAddr := getEnv(REDIS_ADDR)

	return &Config{
		MusicLibraryPath: musicPath,
		RedisAddr:        redisAddr,
	}
}
