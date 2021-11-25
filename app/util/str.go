package util

import (
	"context"
	"math/rand"
	"strconv"
	"time"
	"ws/app/databases"
)

func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func GetSystemReqId() string {
	key := "system:req-id"
	ctx := context.Background()
	cmd := databases.Redis.Incr(ctx, key)
	return "s" + strconv.FormatInt(cmd.Val(), 10)
}