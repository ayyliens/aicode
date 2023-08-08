package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
	"github.com/mitranim/gt"
)

type MockClient func(ChatCompletionMessage) ChatCompletionMessage

var _ Client = gg.Zero[MockClient]()

func (self MockClient) ChatCompletionResponse(ctx u.Ctx, src ChatCompletionRequest) (out ChatCompletionResponse) {
	out.Id = gt.RandomUuid().String()

	gg.Append(&out.Choices, ChatCompletionChoice{
		Index:        0,
		FinishReason: FinishReasonStop,
		Message:      gg.Ptr(self.ChatCompletionMessage(gg.Last(src.Messages))),
	})

	return
}

func (self MockClient) ChatCompletionBody(ctx u.Ctx, src ChatCompletionRequest) []byte {
	return gg.JsonBytes(self.ChatCompletionResponse(ctx, src))
}

func (self MockClient) ChatCompletionMessage(src ChatCompletionMessage) (out ChatCompletionMessage) {
	if self != nil {
		return self(src)
	}

	out.Role = ChatMessageRoleAssistant
	out.Content = `(mock response)`
	return
}
