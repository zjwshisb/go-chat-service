package cron

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
	"ws/app/http/websocket"
)

func clearChannel() {
	websocket.AdminManager.ClearInactiveChannel()
	websocket.UserManager.ClearInactiveChannel()
}

func Serve() *gocron.Scheduler {
	fmt.Println("start cron")
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(closeSessions)
	s.Every(1).Minute().Do(clearChannel)
	s.StartAsync()
	return s
}
