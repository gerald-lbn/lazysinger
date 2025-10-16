package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	Run(ctx context.Context)
	AddJob(crontab string, cmd func()) (int, error)
	RemoveJob(id int)
}

func NewScheduler() Scheduler {
	return &scheduler{
		cron: cron.New(cron.WithSeconds()),
	}
}

type scheduler struct {
	cron *cron.Cron
}

func (s *scheduler) Run(ctx context.Context) {
	s.cron.Start()
	<-ctx.Done()
	s.cron.Stop()
}

func (s *scheduler) AddJob(crontab string, cmd func()) (int, error) {
	id, err := s.cron.AddFunc(crontab, cmd)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *scheduler) RemoveJob(id int) {
	s.cron.Remove(cron.EntryID(id))
}
