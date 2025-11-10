package services

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

const (
	DownloadLyricsQueueName = "download_lyrics"
)

// WorkerService manages workers
type WorkerService struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewWorkerService() (*WorkerService, error) {
	mux := asynq.NewServeMux()

	return &WorkerService{
		server: nil,
		mux:    mux,
	}, nil
}

func (ws *WorkerService) Start() {
	go ws.server.Run(ws.mux)
	slog.Info("worker server started")
}
