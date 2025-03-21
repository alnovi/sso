package scheduler

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Task interface {
	Handle() error
}

type Scheduler struct {
	cron gocron.Scheduler
}

func New(stopTimeout time.Duration) (*Scheduler, error) {
	cron, err := gocron.NewScheduler(
		gocron.WithStopTimeout(stopTimeout),
	)
	return &Scheduler{cron: cron}, err
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() error {
	return s.cron.Shutdown()
}

func (s *Scheduler) AddCronTask(crontab string, task Task) error {
	_, err := s.cron.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(task.Handle),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	return err
}

func (s *Scheduler) AddDurationTask(duration time.Duration, task Task) error {
	_, err := s.cron.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(task.Handle),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	return err
}

func (s *Scheduler) AddDailyTask(interval, hours, minutes uint, task Task) error {
	_, err := s.cron.NewJob(
		gocron.DailyJob(interval, gocron.NewAtTimes(gocron.NewAtTime(hours, minutes, 0))),
		gocron.NewTask(task.Handle),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	return err
}
