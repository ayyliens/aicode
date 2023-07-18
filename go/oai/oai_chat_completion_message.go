package oai

import "github.com/mitranim/gg"

type ChatCompletionMessage struct {
	Role         ChatMessageRole `json:"role,omitempty"          yaml:"role,omitempty"          toml:"role,omitempty"`
	Content      string          `json:"content,omitempty"       yaml:"content,omitempty"       toml:"content,omitempty"`
	FunctionCall *FunctionCall   `json:"function_call,omitempty" yaml:"function_call,omitempty" toml:"function_call,omitempty"`
}

func (self ChatCompletionMessage) IsValid() bool {
	defer gg.Skip()
	self.Validate()
	return true
}

func (self ChatCompletionMessage) Validate() {
	self.ValidateRole()
	self.ValidateContent()
}

func (self ChatCompletionMessage) ValidateRole() {
	if gg.IsZero(self.Role) {
		panic(gg.Errf(`invalid %T: missing role`, self))
	}
}

func (self ChatCompletionMessage) ValidateContent() {
	if gg.IsZero(self.Content) &&
		(self.FunctionCall == nil || !self.FunctionCall.IsValid()) {
		panic(gg.Errf(`invalid %T: missing content and function call`, self))
	}
}

func (self ChatCompletionMessage) Ext() string {
	if self.FunctionCall == nil {
		return `.md`
	}
	return `.yaml`
}
