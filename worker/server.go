package worker

import (
	"github.com/gerald-lbn/lazysinger/config"
	"github.com/hibiken/asynq"
)

func RunServer() error {
	asynqServer := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: config.Server.Worker.RedisAddr,
		},
		asynq.Config{
			Concurrency: config.Server.Worker.Concurrency,
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	// mux handlers
	mux.HandleFunc(TypeLyricsDownload, HandleDownloadLyricsTask)

	return asynqServer.Run(mux)
}
