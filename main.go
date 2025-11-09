package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/refrain/pkg/watcher"
	"github.com/gerald-lbn/refrain/pkg/worker"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := mainContext(context.Background())
	defer cancel()

	runRefrain(ctx)
}

// mainContext returns a context that is cancelled when the process receives a signal to exit.
func mainContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)
}

func runRefrain(ctx context.Context) {
	log.Print("Starting Refrain")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(startWatcher(ctx))
	g.Go(startWorkerServer(ctx))

	if err := g.Wait(); err != nil {
		log.Print("A fatal error occured in Refrain. Aborting", err)
	}
}

func startWatcher(ctx context.Context) func() error {
	log.Print("Starting watcher")

	return func() error {
		w, err := watcher.NewWatcher()
		if err != nil {
			return err
		}

		path := "/music"
		if err := w.Start(ctx, path); err != nil {
			log.Fatalf("An error occured while watching: '%s'. Reason: %v", path, err)
		}
		defer w.Stop()

		return nil
	}
}

func startWorkerServer(ctx context.Context) func() error {
	server := worker.NewServer()
	mux := worker.NewServeMux()

	mux.HandleFunc(worker.TypeDownloadLyrics, worker.HandleDownloadLyricsTask)

	return worker.StartAsynqServer(ctx, server, mux)
}
