package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
	"github.com/mitranim/gr"
)

type OaiClient struct{ ApiKey string }

func (self OaiClient) ValidApiKey() string {
	if gg.IsZero(self.ApiKey) {
		panic(gg.Errf(`missing API key in %T`, self))
	}
	return self.ApiKey
}

func (self OaiClient) Req() *gr.Req {
	return new(gr.Req).
		To(`https://api.openai.com/v1`).
		HeadSet(`Authorization`, `Bearer `+self.ValidApiKey()).
		HeadSet(`Accept`, gr.TypeJsonUtf8)
}

func (self OaiClient) ChatCompletionReq() *gr.Req {
	return self.Req().Join(`/chat/completions`).Post()
}

// Caller must close.
func (self OaiClient) ChatCompletionRes(ctx u.Ctx, src ChatCompletionRequest) *gr.Res {
	return self.ChatCompletionReq().Ctx(ctx).Json(src).Res().Ok()
}

func (self OaiClient) ChatCompletionBody(ctx u.Ctx, src ChatCompletionRequest) []byte {
	return self.ChatCompletionRes(ctx, src).ReadBytes()
}

func (self OaiClient) ChatCompletion(ctx u.Ctx, src ChatCompletionRequest) (out ChatCompletionResponse) {
	self.ChatCompletionRes(ctx, src).Json(&out)
	return
}
