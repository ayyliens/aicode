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

  - Support simultaneously watching and operating on multiple directories.
    The user should watch an "ancestor" dir, and we should operate on all
    nested directories that appear to be "conv" dirs. Motives:

  - Bot can be slow. Operator may work on multiple unrelated directories at
    once to minimize waiting.

  - May allow easy switching between different forks of the same directory.
*/
type ClientConvDir struct {
	ClientCommon
	// FIXME trunc only assistant files
	Trunc              bool              `flag:"--trunc"                desc:"support conversation truncation in watch mode (best with --fork)" json:"trunc,omitempty"              yaml:"trunc,omitempty"              toml:"trunc,omitempty"`
	Fork               bool              `flag:"--fork"                 desc:"support conversation forking in watch mode (best with --trunc)"   json:"fork,omitempty"               yaml:"fork,omitempty"               toml:"fork,omitempty"`
	TruncAfter         gg.Opt[u.Version] `flag:"--trunc-after"          desc:"always truncate files after given index (best with --fork)"       json:"truncAfter,omitempty"         yaml:"truncAfter,omitempty"         toml:"truncAfter,omitempty"`
	Dry                bool              `flag:"--dry"                  desc:"dry run: no request to external API (cancels --rec)"              json:"dry,omitempty"                yaml:"dry,omitempty"                toml:"dry,omitempty"`
	ReadResponseLatest bool              `flag:"--read-response-latest" desc:"instead of actual HTTP request, read last response from disk"     json:"readResponseLatest,omitempty" yaml:"readResponseLatest,omitempty" toml:"readResponseLatest,omitempty"`
	Rec                bool              `flag:"--rec"                  desc:"recursive run: execute prompts until skip reason"                 json:"rec,omitempty"                yaml:"rec,omitempty"                toml:"rec,omitempty"`
	Tail               gg.Opt[int]       `flag:"--tail"                 desc:"limit of messages to include into request context"                json:"tail,omitempty"               yaml:"tail,omitempty"               toml:"tail,omitempty"`
	Functions          Functions         `json:"-" yaml:"-" toml:"-"`
}

func (self ClientConvDir) Run(ctx u.Ctx) {
	if self.Rec && self.Dry {
		log.Println(`dry run: rec is ignored`)
	}

	if self.Watch {
		self.RunWatch(ctx)
	} else {
		self.RunOnce(ctx)
	}
}

func (self ClientConvDir) RunWatch(ctx u.Ctx) {
	self.InitBackupOpt()
	self.InitNextMessagePlaceholderOpt()

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
	dir := self.ConvDir()

	defer gg.Finally(dir.WriteStatusFinally)
	defer gg.Finally(u.LogErr)
	dir.WriteStatusPending()

	if !dir.HasVersionedFiles() {
		if self.Verb {
			log.Println(`no messages found, creating placeholder`)
		}
		dir.WriteNextMessagePlaceholder()
		return
	}

	truncIndex, trunc := self.ShouldTruncAfter(eve)
	if trunc && dir.CanTruncAfter(truncIndex) {
		if self.Fork {
			self.ForkFromBackup(dir.PathToFork())
		}
		if self.Trunc {
			dir.TruncAfter(truncIndex)
		}
	}

	// load files and check validity of messages
	dir.InitFiles(self.Functions)

	ver, msg := dir.PendingMessage()
	skip := msg.SkipReason()
	if gg.IsNotZero(skip) {
		if self.Verb {
			log.Println(`skipping: last message:`, skip)
		}
		return
	}

	if self.Rec && !self.Dry {
		// execute until skip reason
		defer self.RunOnFsEvent(ctx, nil)
	}

	// Somewhat redundant with below, TODO dedup.
	call := msg.GetFunctionCall()
	if gg.IsNotZero(call) {
		self.RunFunction(ctx, ver.AddMinor(), dir, call)
		return
	}

	// generation of request file for further sending
	req := dir.ChatCompletionRequest(ver, self.Model)

	if !gg.IsZero(self.Tail) {
		// TODO Don't include unnecessary messages on `dir.ChatCompletionRequest`.
		// Tentative. Include only last N messages.
		req.Messages = u.TakeLast(req.Messages, self.Tail.Val)
	}

	for _, function := range self.Functions.Slice {
		gg.Append(&req.Functions, function.Def())
	}

	dir.WriteRequestLatest(req)

	if self.Dry {
		if self.Verb {
			log.Println(`dry run: skipping request`)
		}
		return
	}

	resBody := self.VerbChatCompletionBody(ctx, req, dir)
	dir.WriteResponseJson(resBody)

	res := gg.JsonDecodeTo[ChatCompletionResponse](resBody)
	dir.WriteResponseEncoded(res)

	choice := res.ChatCompletionChoice()
	choice.FinishReason.Validate()

	msg = choice.ChatCompletionMessage()
	msg.Validate()

	nextVer := ver.AddMinor()
	dir.WriteMessage(nextVer, msg)
	self.RunFunctionOpt(ctx, nextVer.AddMinor(), dir, msg.GetFunctionCall())
}

func (self ClientConvDir) RunFunctionOpt(ctx u.Ctx, ver u.Version, dir ConvDir, call FunctionCall) {
	if gg.IsZero(call) {
		if self.Verb {
			log.Println(`function call is empty, creating placeholder`)
		}
		dir.InitNextMessagePlaceholder()
		return
	}
	self.RunFunction(ctx, ver, dir, call)
}

func (self ClientConvDir) RunFunction(ctx u.Ctx, ver u.Version, dir ConvDir, call FunctionCall) {
	/**
	If we fail to process the function call, then in addition to logging the
	error, which is done by the caller outside of this function, we also create
	a regular msg placeholder (text/markdown), so the user can continue the
	conversation more easily. This might be part of a normal workflow: bots may
	first produce malformed outputs, and then be cajoled into producing
	something usable.
	*/
	defer u.Fail0(dir.WriteNextMessagePlaceholderOrSkip)

	out := self.Functions.Response(ctx, call.Name, call.Arguments.String(), self.Verbose)

	if gg.IsZero(out) {
		if self.Verb {
			log.Printf(
				`empty response from function %q; avoiding automatic function response to avoid confusing the bot; preferring manual text response`,
				call.Name,
			)
		}
		dir.InitNextMessagePlaceholder()
		return
	}

	dir.WriteMessageFunctionResponse(ver, call.Name, out)

	if self.Verb {
		log.Printf(`wrote pending function response; when running in watch mode, review and re-save the file to trigger the next request`)
	}
}

func (self ClientConvDir) VerbChatCompletionBody(ctx u.Ctx, req ChatCompletionRequest, dir ConvDir) []byte {
	defer gg.RecWith(u.LogErr)

	if self.ReadResponseLatest {
		if self.Verb {
			log.Printf(`reusing last response from disk (flag "--read-response-latest"), path %q`, dir.PathToResponseLatestJson())
		}
		return dir.ReadResponseJson()
	}
	if self.Verb {
		defer gg.LogTimeNow(`chat completion request`).LogStart().LogEnd()
	}
	return self.ChatCompletionBody(ctx, req)
}

func (self ClientConvDir) ConvDir() (out ConvDir) {
	out.Path = self.Path
	out.Verbose = self.Verbose
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

	u.RemoveFileOrDirOrSkip(tar)
	u.CopyDirRec(src, tar)
}

func (self ClientConvDir) ForkFromBackup(tar string) {
	src := self.BackupDirPath()
	if self.Verb {
		log.Printf(`forking directory %q to %q from backup %q`, self.Path, tar, src)
	}
	u.CopyDirRec(src, tar)
}

/*
This operation should be done only in watch mode without `.Init`. When flag
`.Init` is set, we perform this on the initial run, which also has other
potential side effects.
*/
func (self ClientConvDir) InitNextMessagePlaceholderOpt() {
	if !self.Init {
		dir := self.ConvDir()
		dir.InitNextMessagePlaceholder()
	}
}

// Implement `u.FsEventSkipper`.
func (self ClientConvDir) ShouldSkipFsEvent(eve notify.EventInfo) bool {
	if eve == nil {
		return true
	}

	name := filepath.Base(eve.Path())

	return !IsVersionedFileNameLax(name) &&
		u.BaseNameWithoutExt(name) != BaseNameRequestTemplate
}

func (self ClientConvDir) ShouldTruncAfter(eve notify.EventInfo) (_ u.Version, _ bool) {
	if self.TruncAfter.IsNotNull() {
		return self.TruncAfter.Val, true
	}

	if eve != nil && eve.Event() == notify.Write {
		name := ParseIndexedFileNameOpt(filepath.Base(u.NotifyEventPath(eve)))
		if gg.IsNotZero(name) {
			return name.Version, true
		}
	}

	return
}
