package oai

type ChatMessageRole string

const (
	ChatMessageRoleSystem    = `system`
	ChatMessageRoleUser      = `user`
	ChatMessageRoleAssistant = `assistant`
	ChatMessageRoleFunction  = `function`
)
