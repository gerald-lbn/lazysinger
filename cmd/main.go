package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/lazysinger/internal/config"
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

	g.Go(startWatcher())
	g.Go(startWorkerServer())

	if err := g.Wait(); err != nil {
		log.Print("Fatal error in LazySinger. Aborting", err)
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
			log.Fatalf("error creating scanner: %v", err)
		}
		defer watcher.Close()

		watcher.HandleCreate = func(pathToFile string) error {
			log.Printf("New file created: %s. Pushing it to queue", pathToFile)
			task, err := queue.NewDownloadLyricsTask(pathToFile)
			if err != nil {
				return err
			}
			info, err := asynqClient.Enqueue(task)
			if err != nil {
				return err
			}
			log.Printf("Job #%s created and pushed to queue '%s' to process %s", info.ID, info.Queue, pathToFile)
			return nil
		}

		watcher.HandleDelete = func(pathToFile string) error {
			log.Printf("File deleted: %s. Pushing it to queue", pathToFile)
			task, err := queue.NewDeleteLyricsTask(pathToFile)
			if err != nil {
				return err
			}
			info, err := asynqClient.Enqueue(task, asynq.Queue(queue.CRITICAL))
			log.Printf("Job #%s created and pushed to queue '%s' to process %s", info.ID, info.Queue, pathToFile)
			return nil
		}

		watcher.HandleMove = func(pathToFile string) error {
			log.Printf("Filed moved: %s.", pathToFile)
			return nil
		}

		watcher.HandleRename = func(pathToFile string) error {
			log.Printf("Filed renamed: %s", pathToFile)
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
			log.Printf("Job #%s created and pushed to queue '%s' to process %s", info.ID, info.Queue, pathToFile)
			return nil
		}

		if err := watcher.InitialScan(); err != nil {
			log.Fatalf("Initial scan failed: %v", err)
		} else {
			log.Printf("Initial scan completed")
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
