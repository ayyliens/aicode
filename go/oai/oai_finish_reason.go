package oai

type FinishReason string

const (
	FinishReasonNone          FinishReason = ``
	FinishReasonStop          FinishReason = `stop`
	FinishReasonLength        FinishReason = `length`
	FinishReasonFunctionCall  FinishReason = `function_call`
	FinishReasonContentFilter FinishReason = `content_filter`
	FinishReasonNull          FinishReason = `null`
)
