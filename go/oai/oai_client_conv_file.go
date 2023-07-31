package oai

import (
	"_/go/u"
	"strings"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/rjeczalik/notify"
)

// Short for "OpenAI client for/with conversation file".
type ClientConvFile struct{ ClientCommon }

func (self ClientConvFile) Run(ctx u.Ctx) {
	if self.Watch {
		self.RunWatch(ctx)
	} else {
		self.RunOnce(ctx)
	}
}

func (self ClientConvFile) RunWatch(ctx u.Ctx) {
	var wat u.Watcher[ClientConvFile]
	wat.Runner = self
	wat.WatcherCommon = self.WatcherCommon
	wat.Create = true
	wat.Run(ctx)
}

func (self ClientConvFile) OnFsEvent(ctx u.Ctx, _ notify.EventInfo) {
	defer gg.RecWith(u.LogErr)
	self.RunOnce(ctx)
}

func (self ClientConvFile) RunOnce(ctx u.Ctx) {
	src := strings.TrimSpace(gg.ReadFile[string](self.Path))

	var req ChatCompletionRequest
	req.Default()
	req.Messages = u.ParseJsonLines[ChatCompletionMessage](src)

	if gg.IsEmpty(req.Messages) {
		return
	}

	gg.WriteFile(self.Path, gg.JoinLinesOpt(
		src,
		`Sending chat completion request...`,
	))

	defer gg.Fail(func(error) { gg.WriteFile(self.Path, src) })

	res := self.ChatCompletion(ctx, req)
	if self.Verb {
		grepr.Prn(`ChatCompletionResponse:`, res)
	}

	gg.WriteFile(self.Path, gg.JoinLinesOpt(
		src,
		u.JsonEncodePretty[string](res.ChatCompletionMessage()),
	))
}
