package cron

import (
	"github.com/go-co-op/gocron"
	"time"
)

func Run()  {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(closeSessions)
	s.StartAsync()
}
