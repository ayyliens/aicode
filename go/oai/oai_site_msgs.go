package oai

import (
	"_/go/u"
	"os"

	"github.com/mitranim/gg"
)

// Represents the format used internally by the chat UI on the OpenAI website.
type SiteMsgs []SiteMsg

func (self *SiteMsgs) FromFile(file *os.File) {
	u.DecodeJsonLinesInto(file, self)
}

func (self SiteMsgs) ChatCompletionMessages() []ChatCompletionMessage {
	return gg.MapCompact(self, SiteMsg.ChatCompletionMessage)
}

func (self SiteMsgs) ConcatText() string {
	return u.JoinLines2Opt(gg.Map(self, SiteMsg.GetText)...)
}
