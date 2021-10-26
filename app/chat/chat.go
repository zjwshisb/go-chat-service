package chat

import (
	"log"
	"strconv"
)

const (
	// 客服 => {uid: lastTime} hashes
	adminUserLastChatKey = "admin:%d:chat-user:last-time"
	// 待人工接入的用户 sets
	manualUserKey = "user:manual"
)


// 离线时超过多久就自动断开会话
func GetOfflineDuration() int64 {
	setting := Settings[MinuteToBreak]
	minuteStr := setting.GetValue()
	minute, err := strconv.ParseInt(minuteStr, 10,64)
	if err != nil {
		log.Fatal(err)
	}
	return minute * 60
}
