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
	u.Inited
}

func (self OaiClientConvDir) Watch(ctx u.Ctx) {
	self.InitMsg()

	u.Watcher[OaiClientConvDir]{
		Runner: self,
		Path:   self.Path,
		Verb:   self.Verb,
		IsDir:  true,
		Create: true,
		Init:   self.Init,
	}.Run(ctx)
}

func (self OaiClientConvDir) InitMsg() { gg.Ptr(self.ConvDir()).InitMsg() }

func (self OaiClientConvDir) OnFsEvent(ctx u.Ctx, _ notify.EventInfo) {
	defer gg.RecWith(u.LogErr)
	self.Run(ctx)
}

func (self OaiClientConvDir) Run(ctx u.Ctx) {
	defer gg.RecWith(u.LogErr)

	dir := self.ConvDir()
	defer gg.Finally(dir.LogWriteErr)

	req := dir.ChatCompletionRequest()
	skip := req.SkipReason()
	if gg.IsNotZero(skip) {
		if self.Verb {
			log.Println(`skipping chat completion request:`, skip)
		}
		return
	}

	dir.WriteRequestLatest(req)
	res := self.VerbChatCompletionBody(ctx, req)
	dir.WriteResponseLatest(u.JsonPretty(res))
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
