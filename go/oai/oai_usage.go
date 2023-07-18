package oai

type Usage struct {
	PromptTokens     uint64 `json:"prompt_tokens"     yaml:"prompt_tokens"     toml:"prompt_tokens"`
	CompletionTokens uint64 `json:"completion_tokens" yaml:"completion_tokens" toml:"completion_tokens"`
	TotalTokens      uint64 `json:"total_tokens"      yaml:"total_tokens"      toml:"total_tokens"`
}
