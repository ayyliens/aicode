package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
)

type ConvDir struct {
	u.Pathed
	Msgs []ChatCompletionMessage
	Req  gg.Zop[ChatCompletionRequest]
	Res  gg.Zop[ChatCompletionResponse]
}

func (self *ConvDir) Init() { self.Read() }

func (self *ConvDir) Read() {
	self.ReadRequest()
	self.ReadResponse()
	self.ReadMsgs()
}

func (self *ConvDir) ReadRequest() {
	u.PolyDecodeFileOpt(self.RequestPath(), &self.Req.Val)
}

func (self *ConvDir) ReadResponse() {
	u.PolyDecodeFileOpt(self.ResponsePath(), &self.Res.Val)
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

// Can use any extension supported by `u.PolyDecodeFileOpt`.
func (self ConvDir) RequestName() string     { return `request.json` }
func (self ConvDir) RequestNameJson() string { return `request.json` }
func (self ConvDir) RequestPath() string     { return self.PathJoin(self.RequestName()) }
func (self ConvDir) RequestPathJson() string { return self.PathJoin(self.RequestNameJson()) }

// Can use any extension supported by `u.PolyEncodeFileOpt`.
func (self ConvDir) ResponseName() string     { return `response.json` }
func (self ConvDir) ResponseNameJson() string { return `response.json` }
func (self ConvDir) ResponsePath() string     { return self.PathJoin(self.ResponseName()) }
func (self ConvDir) ResponsePathJson() string { return self.PathJoin(self.ResponseNameJson()) }

func (self *ConvDir) InitMsg() {
	if gg.IsEmpty(self.Msgs) {
		self.WriteNextMsgPlaceholderText()
	}
}

func (self ConvDir) ValidMsgs() []ChatCompletionMessage {
	return gg.Filter(self.Msgs, ChatCompletionMessage.IsValid)
}

func (self ConvDir) ChatCompletionRequest() ChatCompletionRequest {
	tar := self.Req.Val
	tar.Default()
	tar.Messages = self.ValidMsgs()
	return tar
}

func (self *ConvDir) WriteResponse(src []byte) {
	res := gg.JsonDecodeTo[ChatCompletionResponse](src)
	self.Res.Set(res)
	self.WriteResponseFiles(src, res)
	self.WriteResponseMsgs(res)
}

func (self *ConvDir) WriteResponseFiles(src []byte, res ChatCompletionResponse) {
	outJson := self.ResponsePathJson()
	u.WriteFile(outJson, src)

	outPoly := self.ResponsePath()
	if outPoly != outJson {
		u.PolyEncodeFileOpt(outPoly, res)
	}
}

func (self *ConvDir) WriteResponseMsgs(res ChatCompletionResponse) {
	choice := res.ChatCompletionChoice()
	msg := choice.ChatCompletionMessage()
	msg.Validate()
	self.WriteNextMsg(msg)

	switch choice.FinishReason {
	case FinishReasonNone, FinishReasonStop:
		self.WriteNextMsgPlaceholderText()

	case FinishReasonFunctionCall:
		self.WriteNextMsgPlaceholderFunctionCall()

	default:
		panic(gg.Errf(`unrecognized/unsupported finish reason %q`, choice.FinishReason))
	}
}

func (self *ConvDir) WriteNextMsg(msg ChatCompletionMessage) {
	var tar MsgFileName
	tar.Index = self.NextIndex()
	tar.Role = msg.Role
	tar.Ext = msg.Ext()

	u.FileWrite{
		Path:  self.PathJoin(tar.String()),
		Body:  gg.ToBytes(msg.Content),
		Mkdir: true,
	}.Run()
	gg.Append(&self.Msgs, msg)
}

func (self *ConvDir) WriteNextMsgPlaceholderText() {
	var msg ChatCompletionMessage
	msg.Role = ChatMessageRoleUser
	self.WriteNextMsg(msg)
}

func (self *ConvDir) WriteNextMsgPlaceholderFunctionCall() {
	var msg ChatCompletionMessage
	msg.Role = ChatMessageRoleUser
	msg.FunctionCall = new(FunctionCall)
	self.WriteNextMsg(msg)
}

func (self ConvDir) NextIndex() int { return len(self.Msgs) }

func (self ConvDir) LogWriteErr(err error) {
	u.LogErr(err)
	defer gg.Skip()
	self.WriteErr(err)
}

func (self ConvDir) WriteErr(err error) {
	u.FileWrite{
		Path:  self.PathJoin(`response_error.txt`),
		Body:  gg.ToBytes(u.FormatVerbose(err)),
		Empty: u.FileWriteEmptyTrunc,
	}.Run()
}
