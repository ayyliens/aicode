package oai

import (
	"_/go/u"
	"os"

	"github.com/mitranim/gg"
)

type OaiSiteMsgs []OaiSiteMsg

func (self *OaiSiteMsgs) FromFile(file *os.File) {
	u.DecodeJsonLinesInto(file, self)
}

func (self OaiSiteMsgs) ChatCompletionMessages() []ChatCompletionMessage {
	return gg.MapCompact(self, OaiSiteMsg.ChatCompletionMessage)
}

func (self OaiSiteMsgs) ConcatText() string {
	return gg.Plus(gg.Map(self, OaiSiteMsg.GetText)...)
}
