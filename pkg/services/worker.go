package services

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

// WorkerService represents a service that performs background tasks.
type WorkerService struct {
	server      *asynq.Server
	mux         *asynq.ServeMux
	client      *asynq.Client
	redisAddr   string
	concurrency int
}

func newServer(redisAddr string, concurrency int) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
		},
	)
}

// NewWorkerService creates a new WorkerService.
func NewWorkerService(redisAddr string, concurrency int) *WorkerService {
	srv := newServer(redisAddr, concurrency)
	mux := asynq.NewServeMux()
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	return &WorkerService{
		server:      srv,
		mux:         mux,
		client:      client,
		redisAddr:   redisAddr,
		concurrency: concurrency,
	}
}

func (ws *WorkerService) RegisterHandler(taskType string, handler asynq.HandlerFunc) {
	ws.mux.HandleFunc(taskType, handler)
}

// Start starts the worker service.
func (ws *WorkerService) Start() error {
	err := ws.server.Start(ws.mux)
	slog.Info("worker service started")
	return err
}

// Stop stops the worker service.
func (ws *WorkerService) Stop() error {
	ws.server.Shutdown()
	return ws.client.Close()
}
