package util

import (
	"math/rand"
	"time"
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

func CreateReqId() int64 {
	rand.Seed(time.Now().UnixNano())
	var min int64 = 10000000000
	var max int64 = 99999999999
	return min + rand.Int63n(max - min)
}