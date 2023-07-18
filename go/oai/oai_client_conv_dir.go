package oai

import (
	"_/go/u"
	"log"

	"github.com/mitranim/gg"
	"github.com/rjeczalik/notify"
)

type OaiClientConvDir struct {
	OaiClient
	u.Pathed
	u.Verbose
}

func (self OaiClientConvDir) Watch(ctx u.Ctx) {
	self.InitMsg()

	u.Watcher[OaiClientConvDir]{
		Runner: self,
		Path:   self.Path,
		Verb:   self.Verb,
		IsDir:  true,
		Create: true,
	}.Run(ctx)
}

func (self OaiClientConvDir) InitMsg() { gg.Ptr(self.ConvDir()).InitMsg() }

func (self OaiClientConvDir) OnFsEvent(ctx u.Ctx, eve notify.EventInfo) {
	defer gg.RecWith(u.LogErr)

	dir := self.ConvDir()
	defer gg.RecWith(dir.LogWriteErr)

	req := dir.ChatCompletionRequest()
	if !req.IsValid() {
		if self.Verb {
			log.Println(`skipping chat completion request`)
		}
		return
	}

	res := self.VerbChatCompletionBody(ctx, req)
	dir.WriteResponse(u.JsonPretty(res))
}

func (self OaiClientConvDir) VerbChatCompletionBody(ctx u.Ctx, req ChatCompletionRequest) []byte {
	if self.Verb {
		defer gg.LogTimeNow(`chat completion request`).LogStart().LogEnd()
	}
	return self.ChatCompletionBody(ctx, req)
}

func (self OaiClientConvDir) ConvDir() (out ConvDir) {
	out.Path = self.Path
	out.Init()
	return
}
