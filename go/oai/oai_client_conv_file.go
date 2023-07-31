package oai

import (
	"_/go/u"
	"strings"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/rjeczalik/notify"
)

// Short for "OpenAI client for/with conversation file".
type OaiClientConvFile struct{ OaiClientCommon }

func (self OaiClientConvFile) Run(ctx u.Ctx) {
	if self.Watch {
		self.RunWatch(ctx)
	} else {
		self.RunOnce(ctx)
	}
}

func (self OaiClientConvFile) RunWatch(ctx u.Ctx) {
	var wat u.Watcher[OaiClientConvFile]
	wat.Runner = self
	wat.WatcherCommon = self.WatcherCommon
	wat.Create = true
	wat.Run(ctx)
}

func (self OaiClientConvFile) OnFsEvent(ctx u.Ctx, _ notify.EventInfo) {
	defer gg.RecWith(u.LogErr)
	self.RunOnce(ctx)
}

func (self OaiClientConvFile) RunOnce(ctx u.Ctx) {
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
