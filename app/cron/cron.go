package cron

import (
	"github.com/go-co-op/gocron"
	"time"
	"ws/app/log"
)

func clearChannel() {

}

func Serve() *gocron.Scheduler {
	log.Log.WithField("a-type", "cron").Info("start")
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(closeSessions)
	s.Every(1).Minute().Do(clearChannel)
	s.StartAsync()
	return s
}
