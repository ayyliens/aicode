package oai

import (
	"_/go/u"
)

// Incomplete definition of an OpenAI client. Expand on demand.
type Client interface {
	ChatCompletionResponse(u.Ctx, ChatCompletionRequest) ChatCompletionResponse
	ChatCompletionBody(u.Ctx, ChatCompletionRequest) []byte
}

// Used by `ClientConvDir` and `ClientConvFile`.
type ClientCommon struct {
	Client
	u.WatcherCommon
	u.Watched
	Model string `flag:"--model" desc:"OpenAI model to use (may be unsupported by some clients)" json:"model,omitempty" yaml:"model,omitempty" toml:"model,omitempty"`
}
