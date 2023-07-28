package oai

import (
	"_/go/u"
	"log"
	"os"
	"path/filepath"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/rjeczalik/notify"
)

/*
Short for "OpenAI client for/with conversation directory".

TODO:

	* Support simultaneously watching and operating on multiple directories.
	  The user should watch an "ancestor" dir, and we should operate on all
	  nested directories that appear to be "conv" dirs.
*/
type OaiClientConvDir struct {
	OaiClient
	u.Pathed
	u.Verbose
	u.Inited
	Functions OaiFunctions
	Trunc     bool // Only for watch mode.
	Fork      bool // Only for watch mode.
	Dry       bool
}

func (self OaiClientConvDir) Watch(ctx u.Ctx) {
	self.InitMessage()
	self.InitBackupOpt()

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
	hasInter := dir.HasIntermediateMessage(baseName)
	isInterWrite := eve != nil && eve.Event() == notify.Write && hasInter

	if isInterWrite {
		if self.Fork {
			self.ForkFromBackup(dir.ForkPath())
		}
		if self.Trunc {
			dir.TruncMessagesAndFilesAfterMessageFileName(baseName, self.Verbose)
		}
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

	req := dir.ChatCompletionRequest(msg)
	dir.WriteRequestLatest(req)

	if self.Dry {
		if self.Verb {
			log.Println(`dry run: skipping request`)
		}
		return
	}

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

	if self.Verb {
		defer gg.LogTimeNow(`running function `, grepr.String(name)).LogStart().LogEnd()
	}
	return fun.OaiCall(arg)
}

func (self OaiClientConvDir) VerbChatCompletionBody(ctx u.Ctx, req ChatCompletionRequest) []byte {
	if self.Verb {
		defer gg.LogTimeNow(`chat completion request`).LogStart().LogEnd()
	}
	return self.ChatCompletionBody(ctx, req)
}

func (self *OaiClientConvDir) OaiConvDirInit() (out OaiConvDir) {
	out.Path = self.Path
	out.Read()
	return
}

func (self OaiClientConvDir) BackupDirPath() string {
	/**
	See `Test_filepath_Join_appending_absolute_path`. On Unix, appending
	an "absolute" path to another path works fine, treating the absolute path
	as if it was relative. However, this may not work on Windows. TODO verify.
	*/
	return filepath.Join(os.TempDir(), TempDirName, gg.Try1(filepath.Abs(self.Path)))
}

func (self OaiClientConvDir) InitBackupOpt() {
	if self.Fork {
		self.InitBackup()
	}
}

func (self OaiClientConvDir) InitBackup() {
	src := self.Path
	tar := self.BackupDirPath()

	if self.Verb {
		log.Printf(`creating backup of directory %q at %q for forking`, src, tar)
	}

	u.RemoveFileOrDirOpt(tar)
	u.CopyDirRec(src, tar)
}

func (self OaiClientConvDir) ForkFromBackup(tar string) {
	if self.Verb {
		log.Printf(`forking directory %q to %q`, self.Path, tar)
	}
	u.CopyDirRec(self.BackupDirPath(), tar)
}
