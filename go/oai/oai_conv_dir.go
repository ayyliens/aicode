package oai

import (
	"_/go/u"
	"log"

	"github.com/mitranim/gg"
)

// Short for "OpenAI conversation directory".
type OaiConvDir struct {
	u.Pathed
	Messages    []ChatCompletionMessage
	ReqTemplate gg.Zop[ChatCompletionRequest]
	ReqLatest   gg.Zop[ChatCompletionRequest]
	ResLatest   gg.Zop[ChatCompletionResponse]
}

func (self *OaiConvDir) Read() {
	self.ReadRequestTemplate()
	self.ReadRequestLatest()
	self.ReadResponseLatest()
	self.ReadMessages()
}

func (self *OaiConvDir) ReadRequestTemplate() {
	tar := &self.ReqTemplate.Val
	u.JsonDecodeFileOpt(self.RequestTemplatePath(`.json`), tar)
	u.YamlDecodeFileOpt(self.RequestTemplatePath(`.yaml`), tar)
	u.TomlDecodeFileOpt(self.RequestTemplatePath(`.toml`), tar)
}

func (self *OaiConvDir) ReadRequestLatest() {
	u.PolyDecodeFileOpt(self.RequestLatestPathJson(), &self.ReqLatest.Val)
}

func (self *OaiConvDir) ReadResponseLatest() {
	u.PolyDecodeFileOpt(self.ResponseLatestPath(), &self.ResLatest.Val)
}

/*
TODO consider preserving file names when reading and writing messages.
File names would be a "secret" field not exposed in JSON. This would be
useful for operations that involve comparing file names, comparing paths,
etc.
*/
func (self *OaiConvDir) ReadMessages() {
	for ind, path := range self.MessageFileNames() {
		self.ReadMessageFile(path).ValidateIndex(ind)
	}
	self.ValidateMessages()
}

func (self *OaiConvDir) ReadMessageFile(name string) (out MessageFileName) {
	gg.Try(out.Parse(name))
	gg.Append(&self.Messages, out.ChatCompletionMessage(self.PathJoin(name)))
	return
}

func (self OaiConvDir) MessageFileNames() []string {
	return gg.Filter(u.ReadDirFileNames(self.Path), IsMessageFileNameLax)
}

/*
Note: the last message is meant to be a placeholder for the user, and is allowed
to have empty content, so we don't validate it.
*/
func (self OaiConvDir) ValidateMessages() {
	for _, msg := range gg.Init(self.Messages) {
		msg.Validate()
	}
}

func (self OaiConvDir) RequestTemplatePath(ext string) string {
	return self.PathJoin(`request_template` + ext)
}

func (self OaiConvDir) RequestLatestPathJson() string {
	return self.PathJoin(`request_latest.json`)
}

// Can change to any extension supported by `u.PolyEncodeFileOpt`.
func (self OaiConvDir) ResponseLatestPath() string {
	return self.PathJoin(`response_latest.json`)
}

func (self OaiConvDir) ResponseLatestPathJson() string {
	return self.PathJoin(`response_latest.json`)
}

func (self OaiConvDir) ErrorPath() string { return self.PathJoin(`error.txt`) }

func (self OaiConvDir) ForkPath() string { return u.IndexedDirForkPath(self.Path) }

func (self *OaiConvDir) InitMessage() {
	if gg.IsEmpty(self.Messages) {
		self.WriteNextMessagePlaceholder()
	}
}

func (self OaiConvDir) ValidMessages() []ChatCompletionMessage {
	return gg.Filter(self.Messages, ChatCompletionMessage.IsValid)
}

func (self OaiConvDir) ChatCompletionRequest() ChatCompletionRequest {
	tar := self.ReqTemplate.Val
	tar.Default()
	tar.Messages = self.ValidMessages()
	return tar
}

func (self OaiConvDir) WriteRequestLatest(src ChatCompletionRequest) {
	u.JsonEncodeFile(self.RequestLatestPathJson(), src)
}

func (self *OaiConvDir) WriteResponseJson(src []byte) {
	u.WriteFile(self.ResponseLatestPathJson(), u.JsonPretty(src))
}

func (self *OaiConvDir) WriteResponseEncoded(res ChatCompletionResponse) {
	out := self.ResponseLatestPath()

	// Assumes that `OaiConvDir.WriteResponseJson` is called earlier.
	// We don't want to overwrite original response JSON with JSON
	// generated by decoding and then encoding again. The original
	// has more information, such as fields not listed in our types.
	if out != self.ResponseLatestPathJson() {
		u.PolyEncodeFileOpt(out, res)
	}
}

// Intended for error paths.
func (self *OaiConvDir) WriteNextMessagePlaceholderOrSkip() {
	defer gg.Skip()
	self.WriteNextMessagePlaceholder()
}

func (self *OaiConvDir) WriteNextMessagePlaceholder() {
	var tar ChatCompletionMessage
	tar.Role = ChatMessageRoleUser
	self.WriteNextMessage(tar)
}

func (self *OaiConvDir) WriteNextMessageFunctionResponse(name FunctionName, body string) {
	var tar ChatCompletionMessage
	tar.Role = ChatMessageRoleFunction
	tar.Name = name
	tar.Content = body
	self.WriteNextMessage(tar)
}

func (self *OaiConvDir) WriteNextMessageFunctionResponsePlaceholder(src FunctionCall) {
	self.WriteNextMessageFunctionResponse(src.Name, ``)
}

func (self *OaiConvDir) WriteNextMessage(src ChatCompletionMessage) {
	ext, body := src.ExtBody()

	var tar MessageFileName
	tar.Index = gg.NumConv[uint](self.NextIndex())
	tar.Role = src.Role
	tar.Ext = ext

	name := tar.String()
	src.FileName = name

	u.WriteFileRec(self.PathJoin(name), body)
	gg.Append(&self.Messages, src)
}

func (self OaiConvDir) NextIndex() int { return len(self.Messages) }

func (self OaiConvDir) LogWriteErr(err error) {
	if u.IsErrContextCancel(err) {
		return
	}

	u.LogErr(err)
	defer gg.Skip()
	self.WriteErr(err)
}

func (self OaiConvDir) WriteErr(err error) {
	u.FileWrite{
		Path:  self.ErrorPath(),
		Body:  gg.ToBytes(u.FormatVerbose(err)),
		Empty: u.FileWriteEmptyDelete,
	}.Run()
}

func (self OaiConvDir) HasIntermediateMessage(name string) bool {
	return gg.IsNotZero(name) &&
		gg.Some(gg.Init(self.Messages), func(val ChatCompletionMessage) bool {
			return val.FileName == name
		})
}

func (self *OaiConvDir) TruncMessagesAndFilesAfterMessageFileName(
	name string,
	verb u.Verbose,
) {
	if !self.HasIntermediateMessage(name) {
		return
	}

	if verb.Verb {
		log.Printf(`truncating messages after %q`, name)
	}

	for gg.IsNotEmpty(self.Messages) {
		msg := gg.Last(self.Messages)
		if msg.FileName == name {
			return
		}

		if verb.Verb {
			log.Printf(`removing message %q`, msg.FileName)
		}

		u.RemoveFileOrDir(self.PathJoin(msg.FileName))
		self.Messages = gg.Init(self.Messages)
	}
}
