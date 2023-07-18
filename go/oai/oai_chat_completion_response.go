package oai

import "github.com/mitranim/gg"

type ChatCompletionResponse struct {
	Id      string                 `json:"id,omitempty"      yaml:"id,omitempty"      toml:"id,omitempty"`
	Object  string                 `json:"object,omitempty"  yaml:"object,omitempty"  toml:"object,omitempty"`
	Created int64                  `json:"created,omitempty" yaml:"created,omitempty" toml:"created,omitempty"`
	Model   string                 `json:"model,omitempty"   yaml:"model,omitempty"   toml:"model,omitempty"`
	Choices []ChatCompletionChoice `json:"choices,omitempty" yaml:"choices,omitempty" toml:"choices,omitempty"`
	Usage   Usage                  `json:"usage,omitempty"   yaml:"usage,omitempty"   toml:"usage,omitempty"`
}

func (self ChatCompletionResponse) ChatCompletionChoice() ChatCompletionChoice {
	size := len(self.Choices)

	if size != 1 {
		panic(gg.Errf(
			`expected %T to have exactly 1 %v, found %v`,
			self, gg.Type[ChatCompletionChoice](), size,
		))
	}

	return self.Choices[0]
}

func (self ChatCompletionResponse) ChatCompletionMessage() ChatCompletionMessage {
	return self.ChatCompletionChoice().ChatCompletionMessage()
}
