package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/lazysinger/config"
	"github.com/gerald-lbn/lazysinger/log"
	"github.com/gerald-lbn/lazysinger/scanner"
	"github.com/gerald-lbn/lazysinger/scheduler"
	"github.com/gerald-lbn/lazysinger/worker"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := mainContext(context.Background())
	defer cancel()

	runLazySinger(ctx)
}

func runLazySinger(ctx context.Context) {
	g, _ := errgroup.WithContext(ctx)

	config.Setup()

	log.SetLevelString(config.Server.Logger.Level)

	g.Go(startScheduler(ctx))
	g.Go(startWorkerServer())
	g.Go(schedulePeriodicScan(ctx))

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

func startScheduler(ctx context.Context) func() error {
	return func() error {
		log.Info().Msg("Starting scheduler")
		schedulerInstance := scheduler.GetInstance()
		schedulerInstance.Run(ctx)
		return nil
	}
}

// Starts a worker server to handle background tasks such as downloading lyrics.
func startWorkerServer() func() error {
	return func() error {
		log.Info().Msg("Starting worker server")
		return worker.RunServer()
	}
}

func schedulePeriodicScan(ctx context.Context) func() error {
	return func() error {
		log.Info().Msg("Scheduling periodic scan")
		schedulerInstance := scheduler.GetInstance()

		_, err := schedulerInstance.AddJob(config.Server.Scanner.Schedule, func() {
			err := scanner.ScanAll(ctx)
			if err != nil {
				log.Error().Err(err).Msg("An error occured while scanning directory")
			} else {
				log.Debug().Msg("Directory scanning completed without any errors")
			}
		})
		if err != nil {
			log.Error().Err(err).Msg("An error occured while scheduling periodic task")
		} else {
			log.Debug().Msg("Periodic scan scheduled successfully")
		}

		return err
	}
}
