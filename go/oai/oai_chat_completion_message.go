package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
)

/*
Note: due to issues/limitations of OpenAI JSON API, we have to use `name` with
`,omitempty` but `content` without `,omitempty`.
*/
type ChatCompletionMessage struct {
	// The following fields are part of the OpenAI API.
	Role         ChatMessageRole `json:"role,omitempty"          yaml:"role,omitempty"          toml:"role,omitempty"`
	Name         FunctionName    `json:"name,omitempty"          yaml:"name,omitempty"          toml:"name,omitempty"`
	Content      string          `json:"content"                 yaml:"content,omitempty"       toml:"content,omitempty"`
	FunctionCall *FunctionCall   `json:"function_call,omitempty" yaml:"function_call,omitempty" toml:"function_call,omitempty"`

	// The following fields are used internally by us.
	FileName string `json:"-" yaml:"-" toml:"-"`
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
	if gg.IsNotZero(self.FileName) {
		return u.FileExt(self.FileName)
	}
	if self.HasFunctionSomething() {
		return `.yaml`
	}
	return `.md`
}

/*
TODO better name.
TODO more configurable.
*/
func (self ChatCompletionMessage) ExtBody() (string, []byte) {
	ext := self.Ext()
	if ext == `.yaml` {
		return ext, u.YamlEncode[[]byte](self)
	}
	return ext, gg.ToBytes(self.Content)
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
