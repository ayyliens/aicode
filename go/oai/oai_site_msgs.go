package oai

import (
	"_/go/u"
	"os"

	"github.com/mitranim/gg"
)

// Represents the format used internally by the chat UI on the OpenAI website.
type OaiSiteMsgs []OaiSiteMsg

func (self *OaiSiteMsgs) FromFile(file *os.File) {
	u.DecodeJsonLinesInto(file, self)
}

func (self OaiSiteMsgs) ChatCompletionMessages() []ChatCompletionMessage {
	return gg.MapCompact(self, OaiSiteMsg.ChatCompletionMessage)
}

func (self OaiSiteMsgs) ConcatText() string {
	return u.JoinLines2Opt(gg.Map(self, OaiSiteMsg.GetText)...)
}
