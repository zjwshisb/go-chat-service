package service

import (
	"context"
	"gf-chat/internal/model"
	"github.com/tmc/langchaingo/llms"
)

type (
	ILangchain interface {
		Ask(ctx context.Context, msg string, messages ...llms.MessageContent) (resp string, err error)
		MessageToContent(messages []*model.CustomerChatMessage) []llms.MessageContent
	}
)

var (
	localLangchain ILangchain
)

func Langchain() ILangchain {
	if localLangchain == nil {
		panic("implement not found for interface langchain, forgot register?")
	}
	return localLangchain
}

func RegisterLangchain(i ILangchain) {
	localLangchain = i
}
