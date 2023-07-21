package oai

import "github.com/mitranim/gg"

type FinishReason string

const (
	FinishReasonNone          FinishReason = ``
	FinishReasonStop          FinishReason = `stop`
	FinishReasonLength        FinishReason = `length`
	FinishReasonFunctionCall  FinishReason = `function_call`
	FinishReasonContentFilter FinishReason = `content_filter`
	FinishReasonNull          FinishReason = `null`
)

func (self FinishReason) Validate() {
	switch self {
	case FinishReasonNone, FinishReasonStop, FinishReasonFunctionCall:
	default:
		panic(gg.Errf(`unrecognized/unsupported finish reason %q`, self))
	}
}
