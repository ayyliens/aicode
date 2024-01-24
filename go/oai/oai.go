package oai

import "_/go/u"

const TempDirName = `aicode`

/*
Implemented by types that correspond to "functions" that we provide/expose to
OpenAI bots. If the output string is non-empty, we send the resulting content
back to the bot, as a function response message. Typically, the output should
be JSON-encoded, but we do not define a schema for it, and bots may be able to
understand arbitrary text. As a special case, if the output is empty, we do not
generate a function response message, because an empty function response seems
to confuse bots.
*/
type OaiFunction interface {
	Name() FunctionName
	OaiCall(u.Ctx, string) string
	Def() FunctionDefinition
}
