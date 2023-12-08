package oai

import "github.com/mitranim/gg"

type VersionedFileType string

const (
	/**
	File body must be either `ChatCompletionMessage.Content` as text, or an entire
	`ChatCompletionMessage` encoded in one of the supported data formats.
	*/
	VersionedFileTypeMessage VersionedFileType = `message`

	/**
	File body must be `ChatCompletionRequest` in one of the supported data
	formats.
	*/
	VersionedFileTypeRequest VersionedFileType = `request`

	/**
	File body must be `ConvFileEval` in one of the supported data formats.
	*/
	VersionedFileTypeEval VersionedFileType = `eval`
)

var VersionedFileTypes = []VersionedFileType{
	VersionedFileTypeMessage,
	VersionedFileTypeRequest,
	VersionedFileTypeEval,
}

func (self VersionedFileType) ErrInvalid() error {
	return gg.Errf(`invalid %T %q; valid values: %q`, self, self, VersionedFileTypes)
}
