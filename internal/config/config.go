package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MusicLibraryPath string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	musicPath := os.Getenv("MUSIC_LIBRARY_PATH")
	if musicPath == "" {
		log.Fatalf("MUSIC_LIBRARY_PATH environment variable not set")
	}

	return &Config{
		MusicLibraryPath: musicPath,
	}
}
