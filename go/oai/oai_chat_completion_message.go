package oai

import (
	"_/go/u"
	"fmt"

	"github.com/mitranim/gg"
)

/*
Note: due to issues/limitations of OpenAI JSON API, we have to use `name` with
`,omitempty` but `content` without `,omitempty`.
*/
type ChatCompletionMessage struct {
	Role         ChatMessageRole `json:"role,omitempty"          yaml:"role,omitempty"          toml:"role,omitempty"`
	Name         FunctionName    `json:"name,omitempty"          yaml:"name,omitempty"          toml:"name,omitempty"`
	Content      string          `json:"content"                 yaml:"content,omitempty"       toml:"content,omitempty"`
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

func (self ChatCompletionMessage) IsRoleValid() bool {
	return gg.IsNotZero(self.Role)
}

func (self ChatCompletionMessage) ValidateRole() {
	if !self.IsRoleValid() {
		panic(gg.Errv(self.PrefixInvalid(), `: missing role`))
	}
}

func (self ChatCompletionMessage) IsContentValid() bool {
	return gg.IsNotZero(self.Content) ||
		gg.IsNotZero(self.Name) ||
		(self.FunctionCall != nil && self.FunctionCall.IsValid())
}

func (self ChatCompletionMessage) ValidateContent() {
	if !self.IsContentValid() {
		panic(gg.Errv(self.PrefixInvalid(), `: must contain content, function call, or function response`))
	}
}

func (self ChatCompletionMessage) PrefixInvalid() string {
	return fmt.Sprintf(`invalid %T`, self)
}

/*
Note: this field uses a pointer rather than `gg.Zop` for compatibility with the
YAML and TOML encoders/decoders.
*/
func (self ChatCompletionMessage) GetFunctionCall() FunctionCall {
	return gg.PtrGet(self.FunctionCall)
}

func (self ChatCompletionMessage) HasFunctionCall() bool {
	return gg.IsNotZero(self.GetFunctionCall())
}

func (self ChatCompletionMessage) HasFunctionResponse() bool {
	return gg.IsNotZero(self.Name)
}

func (self ChatCompletionMessage) HasFunctionSomething() bool {
	return self.HasFunctionResponse() || self.HasFunctionCall()
}

// TODO more configurable.
func (self ChatCompletionMessage) Ext() string {
	if self.HasFunctionSomething() {
		return `.yaml`
	}
	return `.md`
}

func (self ChatCompletionMessage) SkipReason() (_ string) {
	if gg.IsZero(self.Role) {
		return `missing role`
	}

	if self.HasFunctionSomething() {
		return
	}

	if self.Role == ChatMessageRoleAssistant {
		return `already from assistant`
	}

	if u.IsTextBlank(self.Content) {
		return `empty content`
	}

	return
}

func (self ChatCompletionMessage) ChatCompletionMessageExt() (out ChatCompletionMessageExt) {
	out.ChatCompletionMessage = self
	return
}
