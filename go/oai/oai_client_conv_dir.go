package oai

import (
	"_/go/u"
	"log"
	"path/filepath"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/rjeczalik/notify"
)

type OaiClientConvDir struct {
	OaiClient
	u.Pathed
	u.Verbose
	u.Inited
	Functions OaiFunctions
	Trunc     bool
}

func (self OaiClientConvDir) Watch(ctx u.Ctx) {
	self.InitMessage()

	u.Watcher[OaiClientConvDir]{
		Runner: self,
		Path:   self.Path,
		Verb:   self.Verb,
		IsDir:  true,
		Create: true,
		Init:   self.Init,
	}.Run(ctx)
}

func (self OaiClientConvDir) InitMessage() {
	gg.Ptr(self.OaiConvDirInit()).InitMessage()
}

func (self OaiClientConvDir) Run(ctx u.Ctx) { self.RunOnFsEvent(ctx, nil) }

func (self OaiClientConvDir) OnFsEvent(ctx u.Ctx, eve notify.EventInfo) {
	defer gg.RecWith(u.LogErr)
	self.RunOnFsEvent(ctx, eve)
}

func (self OaiClientConvDir) RunOnFsEvent(ctx u.Ctx, eve notify.EventInfo) {
	dir := self.OaiConvDirInit()
	defer gg.Finally(dir.LogWriteErr)

	if gg.IsEmpty(dir.Messages) {
		if self.Verb {
			log.Println(`skipping: no messages found`)
		}
		return
	}

	baseName := filepath.Base(u.NotifyEventPath(eve))
	trunc := self.Trunc &&
		eve != nil &&
		eve.Event() == notify.Write &&
		dir.CanTrunc(baseName)

	// TODO: ideally we would truncate AFTER making the request.
	// Also TODO: truncation should "trash" files, making recovery possible.
	// Also TODO: we should also support automatic conv forking as a superior
	// alternative to deleting files.
	if trunc {
		dir.TruncMessagesAndFilesAfterMessageFileName(
			filepath.Base(eve.Path()),
			self.Verbose,
		)
	}

	msg := gg.Last(dir.Messages)
	skip := msg.SkipReason()
	if gg.IsNotZero(skip) {
		if self.Verb {
			log.Println(`skipping: last message:`, skip)
		}
		return
	}

	if msg.HasFunctionCall() {
		self.RunFunction(dir, msg.GetFunctionCall())
		return
	}

	req := dir.ChatCompletionRequest()
	dir.WriteRequestLatest(req)

	resBody := self.VerbChatCompletionBody(ctx, req)
	dir.WriteResponseJson(resBody)

	res := gg.JsonDecodeTo[ChatCompletionResponse](resBody)
	dir.ResLatest.Set(res)
	dir.WriteResponseEncoded(res)

	choice := res.ChatCompletionChoice()
	choice.FinishReason.Validate()

	msg = choice.ChatCompletionMessage()
	msg.Validate()

	dir.WriteNextMessage(msg)

	call := msg.GetFunctionCall()
	if gg.IsZero(call) {
		dir.WriteNextMessagePlaceholder()
		return
	}

	self.RunFunction(dir, call)
}

/*
TODO consider more flexible error handling. In addition to logging an error, we
could also create a regular msg placeholder (text/markdown) instead of a
function response msg placeholder, so the user can continue the conversation
more easily. This might be part of a normal workflow because bots may produce
malformed outputs at first, and then be cajoled into producing something
usable.
*/
func (self OaiClientConvDir) RunFunction(dir OaiConvDir, call FunctionCall) {
	/**
	If we fail to process the function call, then in addition to logging the
	error, which is done by the caller outside of this function, we also create
	a regular msg placeholder (text/markdown), so the user can continue the
	conversation more easily. This might be part of a normal workflow: bots may
	first produce malformed outputs, and then be cajoled into producing
	something usable.
	*/
	defer u.Fail0(dir.WriteNextMessagePlaceholderOrSkip)

	dir.WriteNextMessageFunctionResponse(
		call.Name,
		self.FunctionResponse(self.Functions.Get(call.Name), call.Name, call.Arguments),
	)

	if self.Verb {
		log.Printf(`wrote pending function response; when running in watch mode, review and re-save the file to trigger the next request`)
	}
}

func (self OaiClientConvDir) FunctionResponse(fun OaiFunction, name FunctionName, arg string) (_ string) {
	if fun == nil {
		if self.Verb {
			log.Printf(`found no registered function %q, returning empty function response`, name)
		}
		return
	}

	// Note: each registered func must be a pointer.
	// This is enforced by `OaiFunctions`.
	u.JsonDecodeAny(arg, fun)

	if self.Verb {
		defer gg.LogTimeNow(`running function `, grepr.String(name)).LogStart().LogEnd()
	}

	return fun.OaiCall()
}

func (self OaiClientConvDir) VerbChatCompletionBody(ctx u.Ctx, req ChatCompletionRequest) []byte {
	if self.Verb {
		defer gg.LogTimeNow(`chat completion request`).LogStart().LogEnd()
	}
	return self.ChatCompletionBody(ctx, req)
}

func (self OaiClientConvDir) OaiConvDirInit() (out OaiConvDir) {
	out.Path = self.Path
	out.Init()
	return
}
