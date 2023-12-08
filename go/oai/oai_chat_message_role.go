package oai

import "github.com/mitranim/gg"

type ChatMessageRole string

const (
	ChatMessageRoleFunction  = `function`
	ChatMessageRoleSystem    = `system`
	ChatMessageRoleUser      = `user`
	ChatMessageRoleAssistant = `assistant`
)

func (self ChatMessageRole) Index() uint16 {
	if self == ChatMessageRoleFunction {
		return 0
	} else if self == ChatMessageRoleSystem {
		return 1
	} else if self == ChatMessageRoleUser {
		return 2
	} else if self == ChatMessageRoleAssistant {
		return 3
	}

	panic(gg.Errf(`unknown index for role: %v`, self))
}

func ChatMessageRoleValidateMatch(path string, act, exp ChatMessageRole) {
	if gg.NotEqNotZero(act, exp) {
		panic(gg.Errf(
			`unexpected role mismatch in msg %q: expected be %q or empty, found %q`,
			path, exp, act,
		))
	}
}
