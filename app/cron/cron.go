package cron

import (
	"github.com/go-co-op/gocron"
	"time"
	"ws/app/log"
)

func Serve() *gocron.Scheduler {
	log.Log.WithField("a-type", "cron").Info("start")
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(closeSessions)
	s.StartAsync()
	return s
}
