package cron

import (
	"github.com/go-co-op/gocron"
	"time"
)
var s *gocron.Scheduler

func Run()  {
	s = gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(closeSessions)
	s.StartAsync()
	s.Stop()
}
func Stop()  {
	if s != nil{
		s.Stop()
	}
}
