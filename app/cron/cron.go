package cron

import (
	"time"
)

func Run()  {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for  {
		<-ticker.C
		go closeSessions()
	}
}
