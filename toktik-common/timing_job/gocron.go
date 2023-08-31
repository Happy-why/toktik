package timing_job

import (
	"github.com/go-co-op/gocron"
	"time"
)

// scheduleJob 定时任务 gocron
type scheduleJob struct {
	schedule *gocron.Scheduler
}

func NewSchedule() *scheduleJob {
	return &scheduleJob{gocron.NewScheduler(time.Local)}
}

func StartMinuteJob(job interface{}, frequency int, tag ...string) {
	s := gocron.NewScheduler(time.Local)
	s.Every(frequency).Tag(tag...).Minute().Do(job)
	s.StartAsync()
}
