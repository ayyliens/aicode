package oai

import (
	"_/go/u"

	"github.com/mitranim/gg"
)

/*
Represents a file in `ConvDir` with the type `IndexedFileTypeEval`.

Describes an internal request for generating a message file, by using a named
function that must be registered by the caller/user. The function is executed
by the framework to generate the content of the resulting message. See
`ConvDir.EvalFileOpt` for the implementation.
*/
type ConvFileEval struct {
	FileName     IndexedFileName `json:"-"                       yaml:"-"                       toml:"-"`
	Type         IndexedFileType `json:"type,omitempty"          yaml:"type,omitempty"          toml:"type,omitempty"`
	FunctionCall *FunctionCall   `json:"function_call,omitempty" yaml:"function_call,omitempty" toml:"function_call,omitempty"`
}

func (self ConvFileEval) Validate() {
	self.FileName.Validate()

	if self.FileName.Type != IndexedFileTypeEval {
		panic(gg.Errf(
			`inconsistency in file name %q for %T: expected type %q, found type %q`,
			self.FileName, self, IndexedFileTypeEval, self.FileName.Type,
		))
	}

	// We may support other types in the future.
	if self.Type != IndexedFileTypeMessage {
		panic(gg.Errf(
			`%T %q must specify type %q; found unexpected type %q`,
			self, self.FileName, IndexedFileTypeMessage, self.Type,
		))
	}

	if gg.IsZero(self.FunctionCall) {
		panic(gg.Errf(`missing function call in %T %q`, self, self.FileName))
	}

	if !self.FunctionCall.IsValid() {
		panic(gg.Errf(`invalid function call in %T %q`, self, self.FileName))
	}
}

func (self ConvFileEval) ValidTargetName() IndexedFileName {
	name := self.FileName
	name.Type = self.Type
	name.Validate()
	return name
}

func (self *ConvFileEval) DecodeFrom(name IndexedFileName, body []byte) {
	u.PolyDecode(body, self, name.Ext)
	self.FileName = name
	self.Validate()
}
