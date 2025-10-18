package scheduler

import "github.com/gerald-lbn/lazysinger/log"

type logger struct{}

func (l *logger) Info(msg string, keysAndValues ...any) {
	log.Debug().Interface("extra", keysAndValues).Msgf("Scheduler: %s", msg)
}

func (l *logger) Error(err error, msg string, keysAndValues ...any) {
	log.Error().Err(err).Interface("extra", keysAndValues).Msgf("Scheduler: %s", msg)
}
