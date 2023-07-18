package oai

import "github.com/mitranim/gg"

type ChatCompletionChoice struct {
	Index        uint64                 `json:"index,omitempty"         yaml:"index,omitempty"         toml:"index,omitempty"`
	Message      *ChatCompletionMessage `json:"message,omitempty"       yaml:"message,omitempty"       toml:"message,omitempty"`
	FinishReason FinishReason           `json:"finish_reason,omitempty" yaml:"finish_reason,omitempty" toml:"finish_reason,omitempty"`
}

func (self ChatCompletionChoice) ChatCompletionMessage() ChatCompletionMessage {
	out := gg.PtrGet(self.Message)
	if gg.IsZero(out) {
		panic(gg.Errf(
			`unexpected missing %v in %v with index %v and finish reason %q`,
			gg.Type[ChatCompletionMessage](), self, self.Index, self.FinishReason,
		))
	}
	return out
}
