package oai

import (
	"_/go/u"
	"fmt"

	"github.com/mitranim/gg"
)

// Extended version of `ChatCompletionMessage` used internally by our framework.
type ChatCompletionMessageExt struct {
	ChatCompletionMessage `yaml:",inline"`

	/**
	When this is specified together with `.Name`, the message is considered to be
	an internal request for a function call, which must be executed by the
	framework to generate the message content. Such messages must always have
	`.Role = ChatMessageRoleFunction`.
	*/
	Arguments u.YamlJsonString `json:"arguments,omitempty" yaml:"arguments,omitempty" toml:"arguments,omitempty"`

	FileName IndexedMessageFileName `json:"-" yaml:"-" toml:"-"`

	/**
	Should be used when the current message is the latest. Should be merged into
	the directory-level request template, if one exists.
	*/
	RequestTemplate any `json:"request_template,omitempty" yaml:"request_template,omitempty" toml:"request_template,omitempty"`

	/**
	Specifies additional message to be generated by the framework.
	Caution: if we want to convert this to a slice, we'd need to modify
	the algorithm that "evaluates" the message to generate the next one,
	to handle partially successful eval.
	*/
	NextMessage *ChatCompletionMessageExt `json:"next_message,omitempty" yaml:"next_message,omitempty" toml:"next_message,omitempty"`
}

func (self ChatCompletionMessageExt) ValidChatCompletionMessage() (_ ChatCompletionMessage) {
	if self.IsValid() {
		return self.ChatCompletionMessage
	}
	return
}

// Override for `ChatCompletionMessage.ValidateRole` with more details.
func (self ChatCompletionMessageExt) ValidateRole() {
	if !self.IsRoleValid() {
		panic(gg.Errv(self.PrefixInvalid(), `: missing role`))
	}
}

// Override for `ChatCompletionMessage.ValidateContent` with more details.
func (self ChatCompletionMessageExt) ValidateContent() {
	if !self.IsContentValid() {
		panic(gg.Errv(self.PrefixInvalid(), `: must contain content, function call, or function response`))
	}
}

// Override for `ChatCompletionMessage.PrefixInvalid` with more details.
func (self ChatCompletionMessageExt) PrefixInvalid() string {
	if gg.IsNotZero(self.FileName) {
		return fmt.Sprintf(`invalid %T %q`, self, self.FileName)
	}
	return fmt.Sprintf(`invalid %T`, self)
}

// TODO more configurable.
func (self ChatCompletionMessageExt) Ext() string {
	if gg.IsNotZero(self.FileName) {
		return self.FileName.Ext
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
func (self ChatCompletionMessageExt) ExtBody() (string, []byte) {
	ext := self.Ext()
	if ext == `.yaml` {
		return ext, u.YamlEncode[[]byte](self)
	}
	return ext, gg.ToBytes(self.Content)
}

func (self ChatCompletionMessageExt) HasInternalFunctionCall() bool {
	return self.HasFunctionResponse() && gg.IsNotZero(self.Arguments)
}