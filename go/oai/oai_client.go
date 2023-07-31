package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
	"github.com/mitranim/gr"
)

// Used by `ClientConvDir` and `ClientConvFile`.
type ClientCommon struct {
	Client
	u.WatcherCommon
	u.Watched
}

type Client struct {
	ApiKey string `flag:"--api-key" desc:"OpenAI API key" json:"apiKey,omitempty" yaml:"apiKey,omitempty" toml:"apiKey,omitempty"`
}

func (self Client) ValidApiKey() string {
	if gg.IsZero(self.ApiKey) {
		panic(gg.Errf(`missing API key in %T`, self))
	}
	return self.ApiKey
}

func (self Client) Req() *gr.Req {
	return new(gr.Req).
		To(`https://api.openai.com/v1`).
		HeadSet(`Authorization`, `Bearer `+self.ValidApiKey()).
		HeadSet(`Accept`, gr.TypeJsonUtf8)
}

func (self Client) ChatCompletionReq() *gr.Req {
	return self.Req().Join(`/chat/completions`).Post()
}

// Caller must close.
func (self Client) ChatCompletionRes(ctx u.Ctx, src ChatCompletionRequest) *gr.Res {
	return self.ChatCompletionReq().Ctx(ctx).Json(src).Res().Ok()
}

func (self Client) ChatCompletionBody(ctx u.Ctx, src ChatCompletionRequest) []byte {
	return self.ChatCompletionRes(ctx, src).ReadBytes()
}

func (self Client) ChatCompletion(ctx u.Ctx, src ChatCompletionRequest) (out ChatCompletionResponse) {
	self.ChatCompletionRes(ctx, src).Json(&out)
	return
}
