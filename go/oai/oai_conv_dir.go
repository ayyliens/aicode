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
	for _, name := range self.VersionedFileNames() {
		self.EvalFileOpt(name, funs)
	}
}

func (self ConvDir) EvalFileOpt(srcName VersionedFileName, funs Functions) {
	if !srcName.IsEval() {
		return
	}

	var eval ConvFileEval
	eval.DecodeFrom(srcName, self.ReadVersionedFile(srcName))

	tarName := eval.ValidTargetName()
	tarPath := self.JoinFilePath(tarName)

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
	return !self.HasVersionedFiles() || self.IsAllRequestsAnswered()
}

func (self ConvDir) IsAllRequestsAnswered() bool {
	return self.LastVersionedFileName().Role == ChatMessageRoleAssistant
}

func (self ConvDir) VersionedFileNames() []VersionedFileName {
	names := gg.Sorted(gg.Map(self.VersionedFileNameCandidates(), ParseIndexedFileNameValid))
	ValidateIndexedFileNames(names)
	return names
}

func (self ConvDir) VersionedFileNameCandidates() []string {
	return gg.Filter(u.ReadDirFileNames(self.Path), IsVersionedFileNameLax)
}

func (self ConvDir) JoinFilePath(name VersionedFileName) string {
	return self.PathJoin(name.ValidString())
}

func (self ConvDir) HasVersionedFiles() bool {
	return gg.IsNotEmpty(self.VersionedFileNameCandidates())
}

func (self ConvDir) HasVersionedFile(name VersionedFileName) bool {
	return self.HasFile(name.ValidString())
}

func (self ConvDir) ReadVersionedFile(name VersionedFileName) []byte {
	return self.ReadFile(name.ValidString())
}

func (self ConvDir) WriteVersionedFile(name VersionedFileName, body []byte) {
	fileName := name.ValidString()
	if gg.FileExists(self.PathJoin(fileName)) {
		panic(gg.Errf(`file already exists: %q`, name))
	}
	self.WriteFile(fileName, body)
}

func (self ConvDir) DeleteVersionedFile(name VersionedFileName) {
	self.DeleteFile(name.ValidString())
}

/*
Load request parameters used for each request
*/
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
	return u.VersionedDirForkPath(self.Path)
}

/*
TODO consider: instead of continuing from the last file (using all files),
continue from the last file before the first "hole" in file indexes. Could
be useful for edge cases like pre-creating a conversation template.
*/
func (self ConvDir) ChatCompletionRequest(version u.Version, model string) (out ChatCompletionRequest) {
	out.Model = model
	self.ReadRequestTemplate(&out)
	out.Default()

	var meta RequestMeta
	{
		name := self.RequestForVersion(version)
		if gg.IsNotZero(name) {
			body := self.ReadVersionedFile(name)
			out.DecodeFrom(name, body)

			//count of context messages to grab in current request, e.g. null - all conversation, 0 - similar to new conversation, 1 - grab messages from 1 previous request
			meta.DecodeFrom(name, body)
		}
	}
	// FIXME append to existing messages
	_, out.Messages = self.ValidMessagesForDepth(meta.Depth)
	return
}

/*
Converts "conversation files" into messages by merging each group by index and role,
resulting in one message with same role per group. For example, if there are
two files like this: `0001_user_message.md` and `0001_user_message.yaml`, both
files will be merged into one message, where the content comes from `.md`
as-is, and some other fields may be set from the `.yaml` file, which would be
interpreted as the YAML encoding of the `ChatCompletionMessage` type. We
support some other formats as well.
*/
func (self ConvDir) ValidMessages() (version u.Version, out []ChatCompletionMessage) {
	return self.ValidMessagesForDepth(gg.Opt[uint16]{})
}

func (self ConvDir) ValidMessagesForDepth(depth gg.Opt[uint16]) (version u.Version, out []ChatCompletionMessage) {
	messages := gg.Filter(self.VersionedFileNames(), VersionedFileName.IsMessage)
	// TODO now controls the depth of major version
	limit := version.PrevMajor(depth.Val)

	nextAllowed := false
	var prev VersionedFileName
	for _, message := range GroupToSlices(messages, VersionedFileName.Name) {
		var msg ChatCompletionMessage
		for _, part := range message {
			// if assistant message found allow next ver after it
			if !gg.IsZero(prev) && !prev.Version.Equal(part.Version) {
				// version changed
				if part.Role == ChatMessageRoleAssistant {
					nextAllowed = true
				} else {
					if nextAllowed {
						nextAllowed = false
					} else {
						version = prev.Version
						return
					}
				}
			}
			prev = part
			msg.DecodeFrom(part, self.ReadVersionedFile(part))
		}

		if depth.IsNull() || limit.Less(prev.Version) || limit.Equal(prev.Version) {
			gg.Append(&out, msg)
		}
	}
	version = prev.Version
	return
}

func GroupToSlices[Slice ~[]Val, Key comparable, Val any](src Slice, fun func(Val) Key) [][]Val {
	if fun == nil {
		return nil
	}

	var out [][]Val
	tar := map[Key][]Val{}

	if fun != nil {
		for _, val := range src {
			key := fun(val)
			items := tar[key]
			if items != nil {
				gg.Append(&items, val)
			} else {
				item := append(tar[key], val)
				tar[key] = item
				out = append(out, item)
			}
		}
	}
	return out
}

func (self ConvDir) RequestForVersion(version u.Version) (_ VersionedFileName) {
	return gg.Find(self.VersionedFileNames(), func(name VersionedFileName) bool {
		return name.Version.Equal(version) && name.IsRequest()
	})
}

func (self ConvDir) NextVersion() u.Version {
	src := self.LastVersionedFileName()
	if gg.IsNotZero(src) {
		return src.Version.NextMajor()
	}
	return src.Version
}

func (self ConvDir) LastVersionedFileName() (out VersionedFileName) {
	return gg.Last(self.VersionedFileNames())
}

func (self ConvDir) PendingMessage() (ver u.Version, out ChatCompletionMessage) {
	version, messages := self.ValidMessages()
	return version, gg.Last(messages)
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
	self.InitNextMessagePlaceholder()
}

func (self ConvDir) WriteNextMessagePlaceholder() {
	var msg ChatCompletionMessage
	msg.Role = ChatMessageRoleUser
	self.WriteNextMessage(msg)
}

func (self ConvDir) WriteMessageFunctionResponse(ver u.Version, name FunctionName, body string) {
	var msg ChatCompletionMessage
	msg.Role = ChatMessageRoleFunction
	msg.Name = name
	msg.Content = body
	self.WriteMessage(ver, msg)
}

//func (self ConvDir) WriteNextMessageFunctionResponsePlaceholder(src FunctionCall) {
//	self.WriteMessageFunctionResponse(src.Name, ``)
//}

func (self ConvDir) WriteMessage(index u.Version, msg ChatCompletionMessage) {
	ext, body := msg.ExtBody()

	var name VersionedFileName
	name.Version = index
	name.Role = msg.Role
	name.Type = VersionedFileTypeMessage
	name.Ext = ext

	self.WriteVersionedFile(name, body)
}

func (self ConvDir) WriteNextMessage(msg ChatCompletionMessage) {
	ext, body := msg.ExtBody()

	var name VersionedFileName
	name.Version = self.NextVersion()
	name.Role = msg.Role
	name.Type = VersionedFileTypeMessage
	name.Ext = ext

	self.WriteVersionedFile(name, body)
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

func (self ConvDir) CanTruncAfter(ind u.Version) bool {
	return gg.Some(self.VersionedFileNames(), func(val VersionedFileName) bool {
		return ind.Less(val.Version)
	})
}

func (self ConvDir) TruncAfter(ind u.Version) {
	if self.Verb {
		log.Printf(
			`truncating %q by deleting indexed files after index %v`,
			self.Path, ind,
		)
	}

	for _, truncName := range gg.Reversed(gg.TakeLastWhile(self.VersionedFileNames(), func(val VersionedFileName) bool {
		return ind.Less(val.Version)
	})) {
		if self.Verb {
			log.Printf(`deleting file %q`, truncName)
		}
		self.DeleteVersionedFile(truncName)
	}
}
