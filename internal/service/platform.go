// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IPlatform interface {
		GetPlatform(ctx context.Context) string
	}
)

var (
	localPlatform IPlatform
)

func Platform() IPlatform {
	if localPlatform == nil {
		panic("implement not found for interface IPlatform, forgot register?")
	}
	return localPlatform
}

func RegisterPlatform(i IPlatform) {
	localPlatform = i
}
