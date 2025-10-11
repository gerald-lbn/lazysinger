package main

import (
	"log"

	"github.com/gerald-lbn/lazysinger/internal/config"
	"github.com/gerald-lbn/lazysinger/internal/queue"
	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.LoadConfig()

	asynqServer := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: cfg.RedisAddr,
		},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				queue.CRITICAL: 6,
				queue.DEFAULT:  3,
				queue.LOW:      1,
			},
			StrictPriority: true,
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
