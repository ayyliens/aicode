package oai

import "github.com/mitranim/gg"

const (
	/**
	File body must be either `ChatCompletionMessage.Content` as text, or an entire
	`ChatCompletionMessage` encoded in one of the supported data formats.
	*/
	IndexedFileTypeMessage IndexedFileType = `message`

	/**
	File body must be `ChatCompletionRequest` in one of the supported data
	formats.
	*/
	IndexedFileTypeRequest IndexedFileType = `request`

	/**
	File body must be `ConvFileEval` in one of the supported data formats.
	*/
	IndexedFileTypeEval IndexedFileType = `eval`
)

var IndexedFileTypes = []IndexedFileType{
	IndexedFileTypeMessage,
	IndexedFileTypeRequest,
	IndexedFileTypeEval,
}

type IndexedFileType string

func (self IndexedFileType) ErrInvalid() error {
	return gg.Errf(`invalid %T %q; valid values: %q`, self, self, IndexedFileTypes)
}
