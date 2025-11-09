package worker

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
)

const redisAddr = "redis:6379"

// NewServer returns a new Server.
func NewServer() *asynq.Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return srv
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	return mux
}

// StartAsynqServer runs the Asynq server
func StartAsynqServer(ctx context.Context, srv *asynq.Server, mux *asynq.ServeMux) func() error {
	return func() error {
		log.Print("Starting Asynq server")

		go func() {
			if err := srv.Run(mux); err != nil {
				log.Printf("Asynq server failed to start: %v", err)
			}
		}()

		<-ctx.Done()
		srv.Shutdown()
		return nil
	}
}
