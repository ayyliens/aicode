package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
)

var OpenAiDefaultModel = Model(`gpt-4`)

type ChatCompletionRequest struct {
	Model            Model                   `json:"model,omitempty"             yaml:"model,omitempty"             toml:"model,omitempty"`
	Messages         []ChatCompletionMessage `json:"messages,omitempty"          yaml:"messages,omitempty"          toml:"messages,omitempty"`
	MaxTokens        uint64                  `json:"max_tokens,omitempty"        yaml:"max_tokens,omitempty"        toml:"max_tokens,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"       yaml:"temperature,omitempty"       toml:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"             yaml:"top_p,omitempty"             toml:"top_p,omitempty"`
	N                uint64                  `json:"n,omitempty"                 yaml:"n,omitempty"                 toml:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"            yaml:"stream,omitempty"            toml:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"              yaml:"stop,omitempty"              toml:"stop,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"  yaml:"presence_penalty,omitempty"  toml:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty" yaml:"frequency_penalty,omitempty" toml:"frequency_penalty,omitempty"`
	LogitBias        map[string]int          `json:"logit_bias,omitempty"        yaml:"logit_bias,omitempty"        toml:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"              yaml:"user,omitempty"              toml:"user,omitempty"`
	Functions        []FunctionDefinition    `json:"functions,omitempty"         yaml:"functions,omitempty"         toml:"functions,omitempty"`
	FunctionCall     any                     `json:"function_call,omitempty"     yaml:"function_call,omitempty"     toml:"function_call,omitempty"`
}

func (self *ChatCompletionRequest) Default() {
	self.Model = gg.Or(self.Model, OpenAiDefaultModel)
}

func (self *ChatCompletionRequest) IsValid() bool {
	msg := gg.Last(self.Messages)
	return msg.Role == ChatMessageRoleUser && msg.IsValid()
}

func (self *ChatCompletionRequest) Merge(src any) {
	if src != nil {
		gg.JsonDecode(gg.JsonString(src), self)
	}
}

func (self *ChatCompletionRequest) DecodeFrom(name VersionedFileName, body []byte) {
	u.PolyDecode(body, self, name.Ext)
}

type RequestMeta struct {
	Depth gg.Opt[uint16] `json:"depth,omitempty" yaml:"depth,omitempty" toml:"depth,omitempty"`
}

func (self *RequestMeta) DecodeFrom(name VersionedFileName, body []byte) {
	u.PolyDecode(body, self, name.Ext)
}
