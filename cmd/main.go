package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/lazysinger/internal/config"
	"github.com/gerald-lbn/lazysinger/internal/log"
	"github.com/gerald-lbn/lazysinger/internal/music"
	"github.com/gerald-lbn/lazysinger/internal/queue"
	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := mainContext(context.Background())
	defer cancel()

	runLazySinger(ctx)
}

func runLazySinger(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	cfg := config.LoadConfig()

	log.SetLevelString(cfg.LogLevel)

	g.Go(startWatcher())
	g.Go(startWorkerServer())

	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("Fatal error in LazySinger. Aborting")
	}
}

// mainContext returns a context that is cancelled when the process receives a signal to exit.
func mainContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)
}

// startWatcher starts a file system watcher to monitor changes in the music library directory.
func startWatcher() func() error {
	return func() error {
		cfg := config.LoadConfig()

		asynqClient := asynq.NewClient(asynq.RedisClientOpt{
			Addr: cfg.RedisAddr,
		})
		defer asynqClient.Close()

		watcher, err := music.NewLibraryWatcher(cfg.MusicLibraryPath)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating scanner")
		}
		defer watcher.Close()

		watcher.HandleCreate = func(pathToFile string) error {
			log.Info().Str("path", pathToFile).Msg("New file created. Pushing it to queue")
			task, err := queue.NewDownloadLyricsTask(pathToFile)
			if err != nil {
				return err
			}
			info, err := asynqClient.Enqueue(task)
			if err != nil {
				return err
			}
			log.Info().Str("job_id", info.ID).Str("queue", info.Queue).Str("file", pathToFile).Msg("Job created and pushed to queue")
			return nil
		}

		watcher.HandleMove = func(pathToFile string) error {
			log.Info().Str("path", pathToFile).Msg("File moved. Pushing it to queue to download lyrics")
			return nil
		}

		watcher.HandleRename = func(pathToFile string) error {
			log.Info().Str("path", pathToFile).Msg("File renamed. Pushing it to queue to download lyrics")
			return nil
		}

		watcher.HandleFileOnInitialScan = func(pathToFile string) error {
			task, err := queue.NewDownloadLyricsTask(pathToFile)
			if err != nil {
				return err
			}
			info, err := asynqClient.Enqueue(task)
			if err != nil {
				return err
			}
			log.Info().Str("job_id", info.ID).Str("queue", info.Queue).Str("file", pathToFile).Msg("Job created and pushed to queue")
			return nil
		}

		if err := watcher.InitialScan(); err != nil {
			log.Fatal().Err(err).Msg("Initial scan failed")
		} else {
			log.Debug().Msg("Initial scan successful")
		}

		watcher.Start()
		watcher.Wait()

		return nil
	}
}

// startWorkerServer starts a worker server to handle background tasks such as downloading lyrics.
func startWorkerServer() func() error {
	return func() error {
		cfg := config.LoadConfig()

		asynqServer := asynq.NewServer(
			asynq.RedisClientOpt{
				Addr: cfg.RedisAddr,
			},
			asynq.Config{
				Concurrency: cfg.WorkerConcurrency,
				Queues: map[string]int{
					queue.CRITICAL: 6,
					queue.DEFAULT:  3,
					queue.LOW:      1,
				},
				StrictPriority: cfg.WorkerStrictPriority,
			},
		)

		// mux maps a type to a handler
		mux := asynq.NewServeMux()

		// mux handlers
		mux.HandleFunc(queue.TypeLyricsDownload, queue.HandleDownloadLyrics)

		return asynqServer.Run(mux)
	}
}
