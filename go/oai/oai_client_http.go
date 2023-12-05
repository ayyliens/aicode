package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
	"github.com/mitranim/gr"
)

type HttpClient struct {
	ApiKey string `flag:"--api-key" desc:"OpenAI API key"   json:"apiKey,omitempty" yaml:"apiKey,omitempty" toml:"apiKey,omitempty"`
	Model  string `flag:"--model"   desc:"OpenAI API model" json:"model,omitempty"  yaml:"model,omitempty"  toml:"model,omitempty"`
}

var _ Client = gg.Zero[HttpClient]()

func (self HttpClient) ValidApiKey() string {
	if gg.IsZero(self.ApiKey) {
		panic(gg.Errf(`missing API key in %T`, self))
	}
	return self.ApiKey
}

func (self HttpClient) Req() *gr.Req {
	return new(gr.Req).
		To(`https://api.openai.com/v1`).
		HeadSet(`Authorization`, `Bearer `+self.ValidApiKey()).
		HeadSet(`Accept`, gr.TypeJsonUtf8)
}

func (self HttpClient) ChatCompletionReq() *gr.Req {
	return self.Req().Join(`/chat/completions`).Post()
}

// Required for `Client`.
func (self HttpClient) ChatCompletionResponse(ctx u.Ctx, src ChatCompletionRequest) (out ChatCompletionResponse) {
	self.ChatCompletionRes(ctx, src).Json(&out)
	return
}

// Required for `Client`.
func (self HttpClient) ChatCompletionBody(ctx u.Ctx, src ChatCompletionRequest) []byte {
	return self.ChatCompletionRes(ctx, src).ReadBytes()
}

// Caller must close response.
func (self HttpClient) ChatCompletionRes(ctx u.Ctx, src ChatCompletionRequest) *gr.Res {
	return self.ChatCompletionReq().Ctx(ctx).Json(src).Res().Ok()
}
