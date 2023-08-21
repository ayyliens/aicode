package oai

import "github.com/mitranim/gg"

type ChatMessageRole string

const (
	ChatMessageRoleSystem    = `system`
	ChatMessageRoleUser      = `user`
	ChatMessageRoleAssistant = `assistant`
	ChatMessageRoleFunction  = `function`
)

func ChatMessageRoleValidateMatch(path string, act, exp ChatMessageRole) {
	if gg.NotEqNotZero(act, exp) {
		panic(gg.Errf(
			`unexpected role mismatch in msg %q: expected be %q or empty, found %q`,
			path, exp, act,
		))
	}
}
