package main

import (
	"log"

	"github.com/gerald-lbn/lyrisync/internal/config"
	"github.com/gerald-lbn/lyrisync/internal/queue"
	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.LoadConfig()

	asynqServer := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
		},
		asynq.Config{
			Concurrency: 5,
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	// mux handlers
	mux.HandleFunc(queue.TypeLyricsDownload, queue.HandleDownloadLyrics)

	if err := asynqServer.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
