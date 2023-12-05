package oai

import (
	"_/go/u"
	"log"

	"github.com/mitranim/gg"
)

const (
	BaseNameRequestTemplate = `request_template`
	BaseNameRequestLatest   = `request_latest`
	BaseNameResponseLatest  = `response_latest`
	BaseNameStatusPending   = `status_pending`
	BaseNameStatusError     = `status_error`
	BaseNameStatusDone      = `status_done`
)

/*
Short for "conversation directory". Abstraction with various methods for
operating on a directory containing a conversation with an OpenAI bot.
*/
type ConvDir struct {
	u.Pathed
	u.Verbose
}

/*
Side-effectful initialization.
TODO better name.
*/
func (self ConvDir) InitFiles(funs Functions) {
	self.EvalFiles(funs)
	self.InitNextMessagePlaceholder()
}

func (self ConvDir) EvalFiles(funs Functions) {
	for _, name := range self.IndexedFileNames() {
		self.EvalFileOpt(name, funs)
	}
}

func (self ConvDir) EvalFileOpt(srcName IndexedFileName, funs Functions) {
	if !srcName.IsEval() {
		return
	}

	var eval ConvFileEval
	eval.DecodeFrom(srcName, self.ReadIndexedFile(srcName))

	tarName := eval.ValidTargetName()
	tarPath := self.JoinPathIndexed(tarName)

	if self.HasFile(tarPath) {
		return
	}

	call := eval.FunctionCall

	var msg ChatCompletionMessage
	msg.Role = ChatMessageRoleFunction
	msg.Name = call.Name
	msg.Content = funs.Response(call.Name, call.Arguments.String(), self.Verbose)
	msg.Validate()

	gg.WriteFile(tarPath, u.PolyEncode[[]byte](msg, tarName.Ext))
}

func (self ConvDir) InitNextMessagePlaceholder() {
	if self.NeedNextMessagePlaceholder() {
		self.WriteNextMessagePlaceholder()
	}
}

func (self ConvDir) NeedNextMessagePlaceholder() bool {
	return !self.HasIndexedFiles() || self.IsLastMessageFromAssistant()
}

func (self ConvDir) IsLastMessageFromAssistant() bool {
	return self.LastIndexedFileName().Role == ChatMessageRoleAssistant
}

func (self ConvDir) IndexedFileNames() []IndexedFileName {
	names := gg.Map(self.IndexedFileNameCandidates(), ParseIndexedFileNameValid)
	ValidateIndexedFileNames(names)
	return names
}

func (self ConvDir) IndexedFileNameCandidates() []string {
	return gg.Filter(u.ReadDirFileNames(self.Path), IsIndexedFileNameLax)
}

func (self ConvDir) JoinPathIndexed(name IndexedFileName) string {
	return self.PathJoin(name.ValidString())
}

func (self ConvDir) HasIndexedFiles() bool {
	return gg.IsNotEmpty(self.IndexedFileNameCandidates())
}

func (self ConvDir) HasIndexedFile(name IndexedFileName) bool {
	return self.HasFile(name.ValidString())
}

func (self ConvDir) ReadIndexedFile(name IndexedFileName) []byte {
	return self.ReadFile(name.ValidString())
}

func (self ConvDir) WriteIndexedFile(name IndexedFileName, body []byte) {
	self.WriteFile(name.ValidString(), body)
}

func (self ConvDir) DeleteIndexedFile(name IndexedFileName) {
	self.DeleteFile(name.ValidString())
}

func (self ConvDir) ReadRequestTemplate(out *ChatCompletionRequest) {
	u.JsonDecodeFileOpt(self.PathToRequestTemplate(`.json`), out)
	u.YamlDecodeFileOpt(self.PathToRequestTemplate(`.yaml`), out)
	u.TomlDecodeFileOpt(self.PathToRequestTemplate(`.toml`), out)
}

func (self ConvDir) ReadRequestLatest(out *ChatCompletionRequest) {
	u.PolyDecodeFileOpt(self.PathToRequestLatestJson(), out)
}

func (self ConvDir) ReadResponseLatest(out *ChatCompletionResponse) {
	u.PolyDecodeFileOpt(self.PathToResponseLatest(), out)
}

func (self ConvDir) PathToRequestTemplate(ext string) string {
	return self.PathJoin(BaseNameRequestTemplate + ext)
}

func (self ConvDir) PathToRequestLatestJson() string {
	return self.PathJoin(BaseNameRequestLatest + `.json`)
}

// Can change to any extension supported by `u.PolyEncodeFileOpt`.
func (self ConvDir) PathToResponseLatest() string {
	return self.PathJoin(BaseNameResponseLatest + `.json`)
}

func (self ConvDir) PathToResponseLatestJson() string {
	return self.PathJoin(BaseNameResponseLatest + `.json`)
}

func (self ConvDir) PathToStatusPending() string {
	return self.PathJoin(BaseNameStatusPending + `.txt`)
}

func (self ConvDir) PathToStatusError() string {
	return self.PathJoin(BaseNameStatusError + `.txt`)
}

func (self ConvDir) PathToStatusDone() string {
	return self.PathJoin(BaseNameStatusDone + `.txt`)
}

func (self ConvDir) PathToFork() string {
	return u.IndexedDirForkPath(self.Path)
}

func (self ConvDir) ChatCompletionRequest() (out ChatCompletionRequest) {
	self.InitChatCompletionRequest(&out)
	return
}

/*
TODO consider: instead of continuing from the last file (using all files),
continue from the last file before the first "hole" in file indexes. Could
be useful for edge cases like pre-creating a conversation template.
*/
func (self ConvDir) InitChatCompletionRequest(out *ChatCompletionRequest) {
	if out == nil {
		return
	}

	self.ReadRequestTemplate(out)
	out.Default()
	out.Messages = self.ValidMessages()

	name := self.IndexedFileNameForNextRequest()
	if gg.IsNotZero(name) {
		out.DecodeFrom(name, self.ReadIndexedFile(name))
	}
}

/*
Converts "conversation files" into messages by merging each group of identically
indexed files, resulting in one message per group. For example, if there are
two files like this: `0001_user_message.md` and `0001_user_message.yaml`, both
files will be merged into one message, where the content comes from `.md`
as-is, and some other fields may be set from the `.yaml` file, which would be
interpreted as the YAML encoding of the `ChatCompletionMessage` type. We
support some other formats as well.
*/
func (self ConvDir) ValidMessages() (out []ChatCompletionMessage) {
	var prev IndexedFileName
	var msg ChatCompletionMessage

	for _, next := range self.IndexedFileNames() {
		if !next.IsMessage() {
			continue
		}

		if gg.IsZero(prev) {
			prev = next
			msg.DecodeFrom(next, self.ReadIndexedFile(next))
			continue
		}

		if prev.Index != next.Index {
			gg.Append(&out, msg)
			gg.PtrClear(&msg)
			msg.DecodeFrom(next, self.ReadIndexedFile(next))
			continue
		}

		msg.DecodeFrom(next, self.ReadIndexedFile(next))
	}

	if gg.IsNotZero(msg) {
		gg.Append(&out, msg)
	}
	return
}

func (self ConvDir) IndexedFileNameForNextRequest() (_ IndexedFileName) {
	return gg.Find(self.LastIndexedFileNameGroup(), IndexedFileName.IsRequest)
}

func (self ConvDir) NextIndex() u.FileIndex {
	src := self.LastIndexedFileName()
	if gg.IsNotZero(src) {
		return src.Index + 1
	}
	return src.Index
}

func (self ConvDir) LastIndex() u.FileIndex {
	return self.LastIndexedFileName().Index
}

func (self ConvDir) LastIndexedFileName() (out IndexedFileName) {
	src := gg.Last(self.IndexedFileNameCandidates())
	if gg.IsNotZero(src) {
		gg.Parse(src, &out)
	}
	return
}

func (self ConvDir) LastIndexedFileNameGroup() []IndexedFileName {
	src := self.IndexedFileNames()
	ind := gg.Last(src).Index
	return gg.TakeLastWhile(src, func(val IndexedFileName) bool {
		return val.Index == ind
	})
}

func (self ConvDir) LastMessage() (out ChatCompletionMessage) {
	for _, name := range self.LastIndexedFileNameGroup() {
		if name.IsMessage() {
			out.DecodeFrom(name, self.ReadIndexedFile(name))
		}
	}
	return
}

func (self ConvDir) WriteRequestLatest(src ChatCompletionRequest) {
	u.PolyEncodeFileOpt(self.PathToRequestLatestJson(), src)
}

func (self ConvDir) ReadResponseJson() []byte {
	return gg.ReadFile[[]byte](self.PathToResponseLatestJson())
}

func (self ConvDir) WriteResponseJson(src []byte) {
	u.WriteFile(self.PathToResponseLatestJson(), u.JsonPretty(src))
}

func (self ConvDir) WriteResponseEncoded(src ChatCompletionResponse) {
	out := self.PathToResponseLatest()

	// Assumes that `ConvDir.WriteResponseJson` is called earlier.
	// We don't want to overwrite original response JSON with JSON
	// generated by decoding and then encoding again. The original
	// has more information, such as fields not listed in our types.
	if out != self.PathToResponseLatestJson() {
		u.PolyEncodeFileOpt(out, src)
	}
}

// Intended to be called during panic handling, like via `u.Fail0`.
func (self ConvDir) WriteNextMessagePlaceholderOrSkip() {
	defer gg.Skip()
	self.WriteNextMessagePlaceholder()
}

func (self ConvDir) WriteNextMessagePlaceholder() {
	var msg ChatCompletionMessage
	msg.Role = ChatMessageRoleUser
	self.WriteNextMessage(msg)
}

func (self ConvDir) WriteNextMessageFunctionResponse(name FunctionName, body string) {
	var msg ChatCompletionMessage
	msg.Role = ChatMessageRoleFunction
	msg.Name = name
	msg.Content = body
	self.WriteNextMessage(msg)
}

func (self ConvDir) WriteNextMessageFunctionResponsePlaceholder(src FunctionCall) {
	self.WriteNextMessageFunctionResponse(src.Name, ``)
}

func (self ConvDir) WriteNextMessage(msg ChatCompletionMessage) {
	ext, body := msg.ExtBody()

	var name IndexedFileName
	name.Index = self.NextIndex()
	name.Role = msg.Role
	name.Type = IndexedFileTypeMessage
	name.Ext = ext

	self.WriteIndexedFile(name, body)
}

func (self ConvDir) WriteStatusFinally(err error) {
	if u.IsErrContextCancel(err) {
		return
	}

	if err != nil {
		self.WriteStatusError(err)
		return
	}

	self.WriteStatusDone()
}

/*
TODO: useful logging.
For now, this file merely serves as an indicator of ongoing work.
*/
func (self ConvDir) WriteStatusPending() {
	u.RemoveFileOrDirOrSkip(self.PathToStatusError())
	u.RemoveFileOrDirOrSkip(self.PathToStatusDone())
	u.TouchedFileRec(self.PathToStatusPending())
}

/*
TODO: useful logging.
For now, this file merely serves as an indicator of finished work.
*/
func (self ConvDir) WriteStatusDone() {
	u.RemoveFileOrDirOrSkip(self.PathToStatusPending())
	u.RemoveFileOrDirOrSkip(self.PathToStatusError())
	u.TouchedFileRec(self.PathToStatusDone())
}

func (self ConvDir) WriteStatusError(err error) {
	u.RemoveFileOrDirOrSkip(self.PathToStatusPending())
	u.RemoveFileOrDirOrSkip(self.PathToStatusDone())

	path := self.PathToStatusError()

	if err == nil {
		u.RemoveFileOrDirOrSkip(path)
		return
	}

	u.FileWrite{
		Path:  path,
		Body:  gg.ToBytes(u.FormatVerbose(err)),
		Empty: u.FileWriteEmptyDelete,
	}.Run()
}

func (self ConvDir) CanTruncAfter(ind u.FileIndex) bool {
	return gg.Some(self.IndexedFileNames(), func(val IndexedFileName) bool {
		return val.Index > ind
	})
}

func (self ConvDir) TruncAfter(ind u.FileIndex) {
	if self.Verb {
		log.Printf(
			`truncating %q by deleting indexed files after index %v`,
			self.Path, ind,
		)
	}

	for _, truncName := range gg.Reversed(gg.TakeLastWhile(self.IndexedFileNames(), func(val IndexedFileName) bool {
		return val.Index > ind
	})) {
		if self.Verb {
			log.Printf(`deleting file %q`, truncName)
		}
		self.DeleteIndexedFile(truncName)
	}
}
