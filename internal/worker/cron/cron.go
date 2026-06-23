package cron

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

type Job interface {
	Enabled() bool
	Name() string
	Spec() string
	Run(ctx context.Context) error
}

type Runner struct {
	jobs []Job
}

func NewRunner(jobs ...Job) *Runner {
	return &Runner{jobs: jobs}
}

// Start lên lịch toàn bộ job rồi chạy nền. Dừng khi ctx bị cancel.
func (r *Runner) Start(ctx context.Context) error {
	scheduler := cron.New()

	for _, job := range r.jobs {
		job := job
		if !job.Enabled() {
			log.Printf("[Cron] skipped %q (disabled)", job.Name())
			continue
		}
		if _, err := scheduler.AddFunc(job.Spec(), func() {
			if err := job.Run(ctx); err != nil {
				log.Printf("[Cron] %s: %v", job.Name(), err)
			}
		}); err != nil {
			return fmt.Errorf("schedule cron %q: %w", job.Name(), err)
		}
		log.Printf("[Cron] scheduled %q (%s)", job.Name(), job.Spec())
	}

	scheduler.Start()

	go func() {
		<-ctx.Done()
		scheduler.Stop()
		log.Printf("[Cron] stopped")
	}()
	return nil
}
