package oai

import (
	"_/go/u"
	"strings"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/rjeczalik/notify"
)

type OaiClientConvFile struct {
	OaiClient
	u.Pathed
	u.Verbose
	u.Inited
}

func (self OaiClientConvFile) Watch(ctx u.Ctx) {
	u.Watcher[OaiClientConvFile]{
		Runner: self,
		Path:   self.Path,
		Verb:   self.Verb,
		Create: true,
		Init:   self.Init,
	}.Run(ctx)
}

func (self OaiClientConvFile) OnFsEvent(ctx u.Ctx, _ notify.EventInfo) {
	defer gg.RecWith(u.LogErr)
	self.Run(ctx)
}

func (self OaiClientConvFile) Run(ctx u.Ctx) {
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
