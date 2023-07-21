package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
)

type ConvDir struct {
	u.Pathed
	Msgs        []ChatCompletionMessage
	ReqTemplate gg.Zop[ChatCompletionRequest]
	ReqLatest   gg.Zop[ChatCompletionRequest]
	ResLatest   gg.Zop[ChatCompletionResponse]
}

func (self *ConvDir) Init() { self.Read() }

func (self *ConvDir) Read() {
	self.ReadRequestTemplate()
	self.ReadRequestLatest()
	self.ReadResponseLatest()
	self.ReadMsgs()
}

func (self *ConvDir) ReadRequestTemplate() {
	tar := &self.ReqTemplate.Val
	u.JsonDecodeFileOpt(self.RequestTemplatePath(`.json`), tar)
	u.YamlDecodeFileOpt(self.RequestTemplatePath(`.yaml`), tar)
	u.TomlDecodeFileOpt(self.RequestTemplatePath(`.toml`), tar)
}

func (self *ConvDir) ReadRequestLatest() {
	u.PolyDecodeFileOpt(self.RequestLatestPathJson(), &self.ReqLatest.Val)
}

func (self *ConvDir) ReadResponseLatest() {
	u.PolyDecodeFileOpt(self.ResponseLatestPath(), &self.ResLatest.Val)
}

func (self *ConvDir) ReadMsgs() {
	for ind, path := range self.MsgFileNames() {
		self.ReadMsgFile(path).ValidateIndex(ind)
	}
	self.ValidateMsgs()
}

func (self *ConvDir) ReadMsgFile(name string) (out MsgFileName) {
	gg.Try(out.Parse(name))
	gg.Append(&self.Msgs, out.ChatCompletionMessage(self.PathJoin(name)))
	return
}

func (self ConvDir) MsgFileNames() []string {
	return gg.Filter(u.ReadDirFileNames(self.Path), IsMsgFileNameLax)
}

/*
Note: the last message is meant to be a placeholder for the user, and is allowed
to have empty content, so we don't validate it.
*/
func (self ConvDir) ValidateMsgs() {
	for _, msg := range gg.Init(self.Msgs) {
		msg.Validate()
	}
}

func (self ConvDir) RequestTemplatePath(ext string) string {
	return self.PathJoin(`request_template` + ext)
}

func (self ConvDir) RequestLatestPathJson() string {
	return self.PathJoin(`request_latest.json`)
}

// Can change to any extension supported by `u.PolyEncodeFileOpt`.
func (self ConvDir) ResponseLatestPath() string {
	return self.PathJoin(`response_latest.json`)
}

func (self ConvDir) ResponseLatestPathJson() string {
	return self.PathJoin(`response_latest.json`)
}

func (self ConvDir) ResponseLatestErrorPathJson() string {
	return self.PathJoin(`response_latest_error.json`)
}

func (self *ConvDir) InitMsg() {
	if gg.IsEmpty(self.Msgs) {
		self.WriteNextMsgPlaceholderText()
	}
}

func (self ConvDir) ValidMsgs() []ChatCompletionMessage {
	return gg.Filter(self.Msgs, ChatCompletionMessage.IsValid)
}

func (self ConvDir) ChatCompletionRequest() ChatCompletionRequest {
	tar := self.ReqTemplate.Val
	tar.Default()
	tar.Messages = self.ValidMsgs()
	return tar
}

func (self ConvDir) WriteRequestLatest(src ChatCompletionRequest) {
	u.JsonEncodeFile(self.RequestLatestPathJson(), src)
}

func (self *ConvDir) WriteResponseLatest(src []byte) {
	res := gg.JsonDecodeTo[ChatCompletionResponse](src)
	self.ResLatest.Set(res)
	self.WriteResponseFiles(src, res)
	self.WriteResponseMsgs(res)
}

func (self *ConvDir) WriteResponseFiles(src []byte, res ChatCompletionResponse) {
	outJson := self.ResponseLatestPathJson()
	u.WriteFile(outJson, src)

	outPoly := self.ResponseLatestPath()
	if outPoly != outJson {
		u.PolyEncodeFileOpt(outPoly, res)
	}
}

func (self *ConvDir) WriteResponseMsgs(res ChatCompletionResponse) {
	choice := res.ChatCompletionChoice()
	choice.FinishReason.Validate()

	msg := choice.ChatCompletionMessage()
	msg.Validate()
	self.WriteNextMsg(msg)
}

func (self *ConvDir) WriteNextMsg(msg ChatCompletionMessage) {
	self.WriteMsg(msg)
	self.WriteNextMsgPlaceholder(msg)
}

func (self *ConvDir) WriteNextMsgPlaceholder(src ChatCompletionMessage) {
	call := src.GetFunctionCall()
	if gg.IsZero(call) {
		self.WriteNextMsgPlaceholderText()
	} else {
		self.WriteNextMsgPlaceholderFunctionCall(call)
	}
}

func (self *ConvDir) WriteNextMsgPlaceholderText() {
	var tar ChatCompletionMessage
	tar.Role = ChatMessageRoleUser
	self.WriteMsg(tar)
}

func (self *ConvDir) WriteNextMsgPlaceholderFunctionCall(src FunctionCall) {
	var tar ChatCompletionMessage
	tar.Role = ChatMessageRoleFunction
	tar.Name = src.Name
	self.WriteMsg(tar)
}

func (self *ConvDir) WriteMsg(src ChatCompletionMessage) {
	ext, body := src.ExtBody()

	var tar MsgFileName
	tar.Index = self.NextIndex()
	tar.Role = src.Role
	tar.Ext = ext

	u.FileWrite{
		Path:  self.PathJoin(tar.String()),
		Body:  body,
		Mkdir: true,
	}.Run()
	gg.Append(&self.Msgs, src)
}

func (self ConvDir) NextIndex() int { return len(self.Msgs) }

func (self ConvDir) LogWriteErr(err error) {
	u.LogErr(err)
	defer gg.Skip()
	self.WriteErr(err)
}

func (self ConvDir) WriteErr(err error) {
	u.FileWrite{
		Path:  self.ResponseLatestErrorPathJson(),
		Body:  gg.ToBytes(u.FormatVerbose(err)),
		Empty: u.FileWriteEmptyTrunc,
	}.Run()
}
