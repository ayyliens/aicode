package oai

import (
	"_/go/u"
	"log"
	"os"
	"path/filepath"

	"github.com/mitranim/gg"
	"github.com/rjeczalik/notify"
)

/*
Short for "OpenAI client for/with conversation directory".

TODO:

	* Support simultaneously watching and operating on multiple directories.
		The user should watch an "ancestor" dir, and we should operate on all
		nested directories that appear to be "conv" dirs.
*/
type ClientConvDir struct {
	ClientCommon
	Trunc     bool      `flag:"--trunc" desc:"support conversation truncation in watch mode (best with --fork)" json:"trunc,omitempty" yaml:"trunc,omitempty" toml:"trunc,omitempty"`
	Fork      bool      `flag:"--fork"  desc:"support conversation forking in watch mode (best with --trunc)"   json:"fork,omitempty"  yaml:"fork,omitempty"  toml:"fork,omitempty"`
	Dry       bool      `flag:"--dry"   desc:"dry run: no request to external API"                              json:"dry,omitempty"   yaml:"dry,omitempty"   toml:"dry,omitempty"`
	Functions Functions `json:"-" yaml:"-" toml:"-"`
}

func (self ClientConvDir) Run(ctx u.Ctx) {
	if self.Watch {
		self.RunWatch(ctx)
	} else {
		self.RunOnce(ctx)
	}
}

func (self ClientConvDir) RunWatch(ctx u.Ctx) {
	self.InitBackupOpt()

	var wat u.Watcher[ClientConvDir]
	wat.Runner = self
	wat.WatcherCommon = self.WatcherCommon
	wat.IsDir = true
	wat.Create = true
	wat.Run(ctx)
}

func (self ClientConvDir) RunOnce(ctx u.Ctx) { self.RunOnFsEvent(ctx, nil) }

func (self ClientConvDir) OnFsEvent(ctx u.Ctx, eve notify.EventInfo) {
	defer gg.RecWith(u.LogErr)
	self.RunOnFsEvent(ctx, eve)
}

func (self ClientConvDir) RunOnFsEvent(ctx u.Ctx, eve notify.EventInfo) {
	dir := self.OaiConvDirRead()
	defer gg.Finally(dir.LogWriteErr)

	if gg.IsEmpty(dir.Messages) {
		if self.Verb {
			log.Println(`no messages found, creating placeholder`)
		}
		dir.WriteNextMessagePlaceholder()
		return
	}

	if eve != nil && eve.Event() == notify.Write {
		baseName := IndexedMessageFileNameOpt(filepath.Base(u.NotifyEventPath(eve)))

		if dir.HasIntermediateMessage(baseName) {
			if self.Fork {
				self.ForkFromBackup(dir.ForkPath())
			}
			if self.Trunc {
				dir.TruncMessagesAndFilesAfterIndexedMessageFileName(baseName)
			}
		}
	}

	dir.InitFiles(self.Functions)

	msg := gg.Last(dir.Messages)
	skip := msg.SkipReason()
	if gg.IsNotZero(skip) {
		if self.Verb {
			log.Println(`skipping: last message:`, skip)
		}
		return
	}

	// Somewhat redundant with below, TODO dedup.
	call := msg.GetFunctionCall()
	if gg.IsNotZero(call) {
		self.RunFunction(dir, call)
		return
	}

	req := dir.ChatCompletionRequest()
	req.Merge(msg.RequestTemplate)
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

	msg = choice.ChatCompletionMessage().ChatCompletionMessageExt()
	msg.Validate()

	dir.WriteNextMessage(msg)
	self.RunFunctionOpt(dir, msg.GetFunctionCall())
}

func (self ClientConvDir) RunFunctionOpt(dir ConvDir, call FunctionCall) {
	if gg.IsZero(call) {
		dir.WriteNextMessagePlaceholder()
		return
	}
	self.RunFunction(dir, call)
}

func (self ClientConvDir) RunFunction(dir ConvDir, call FunctionCall) {
	/**
	If we fail to process the function call, then in addition to logging the
	error, which is done by the caller outside of this function, we also create
	a regular msg placeholder (text/markdown), so the user can continue the
	conversation more easily. This might be part of a normal workflow: bots may
	first produce malformed outputs, and then be cajoled into producing
	something usable.
	*/
	defer u.Fail0(dir.WriteNextMessagePlaceholderOrSkip)

	out := self.Functions.Response(call.Name, call.Arguments.String(), self.Verbose)

	if gg.IsZero(out) {
		if self.Verb {
			log.Printf(`function response is empty; preferring manual text response to avoid confusing the bot`)
		}
		dir.WriteNextMessagePlaceholder()
		return
	}

	dir.WriteNextMessageFunctionResponse(call.Name, out)

	if self.Verb {
		log.Printf(`wrote pending function response; when running in watch mode, review and re-save the file to trigger the next request`)
	}
}

func (self ClientConvDir) VerbChatCompletionBody(ctx u.Ctx, req ChatCompletionRequest) []byte {
	if self.Verb {
		defer gg.LogTimeNow(`chat completion request`).LogStart().LogEnd()
	}
	return self.ChatCompletionBody(ctx, req)
}

func (self *ClientConvDir) OaiConvDirRead() (out ConvDir) {
	out.Path = self.Path
	out.Verbose = self.Verbose
	out.Read()
	return
}

func (self ClientConvDir) BackupDirPath() string {
	/**
	See `Test_filepath_Join_appending_absolute_path`. On Unix, appending
	an "absolute" path to another path works fine, treating the absolute path
	as if it was relative. However, this may not work on Windows. TODO verify.
	*/
	return filepath.Join(os.TempDir(), TempDirName, gg.Try1(filepath.Abs(self.Path)))
}

func (self ClientConvDir) InitBackupOpt() {
	if self.Fork {
		self.InitBackup()
	}
}

func (self ClientConvDir) InitBackup() {
	src := self.Path
	tar := self.BackupDirPath()

	if self.Verb {
		log.Printf(`creating backup of directory %q at %q for forking`, src, tar)
	}

	u.RemoveFileOrDirOpt(tar)
	u.CopyDirRec(src, tar)
}

func (self ClientConvDir) ForkFromBackup(tar string) {
	src := self.BackupDirPath()
	if self.Verb {
		log.Printf(`forking directory %q to %q from backup %q`, self.Path, tar, src)
	}
	u.CopyDirRec(src, tar)
}

// Implement `u.FsEventSkipper`.
func (self ClientConvDir) ShouldSkipFsEvent(eve notify.EventInfo) bool {
	if eve == nil {
		return true
	}

	name := filepath.Base(eve.Path())

	return !IsIndexedMessageFileNameLax(name) &&
		u.BaseNameWithoutExt(name) != BaseNameRequestTemplate
}
