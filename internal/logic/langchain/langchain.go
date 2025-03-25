package langchain

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var open bool

type langchian struct {
	llm *ollama.LLM
}

func newLangchian() *langchian {
	config, _ := g.Config().Get(gctx.New(), "langchain.model", "llama3.2")
	llm, err := ollama.New(ollama.WithModel(config.String()), ollama.WithSystemPrompt("你是一个客服系统"))
	if err != nil {
		panic(err)
	}
	return &langchian{llm: llm}
}

func init() {
	service.RegisterLangchain(newLangchian())
}

func (receiver langchian) Ask(ctx context.Context, msg string, messages ...llms.MessageContent) (resp string, err error) {
	messages = append(messages, llms.TextParts(llms.ChatMessageTypeHuman, msg))
	completion, err := receiver.llm.GenerateContent(ctx, messages,
		llms.WithTemperature(0.7),
		llms.WithTopP(0.9))
	if err != nil {
		return
	}
	resp = completion.Choices[0].Content
	return
}
func (receiver langchian) MessageToContent(messages []*model.CustomerChatMessage) []llms.MessageContent {
	return slice.Map(messages, func(index int, item *model.CustomerChatMessage) llms.MessageContent {
		role := llms.ChatMessageTypeHuman
		if item.Source == consts.MessageSourceAi {
			role = llms.ChatMessageTypeAI
		}
		return llms.TextParts(role, item.Content)
	})
}
