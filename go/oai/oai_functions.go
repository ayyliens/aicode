package oai

import (
	"_/go/u"
	"log"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

/*
Registry of "functions" suitable for automatic use by OpenAI bots.
Reference:

	https://platform.openai.com/docs/guides/gpt/function-calling
*/
type Functions gg.OrdMap[FunctionName, OaiFunction]

func (self Functions) Has(key FunctionName) bool {
	return self.OrdMap().Has(key)
}

func (self Functions) Get(key FunctionName) OaiFunction {
	return self.OrdMap().Get(key)
}

func (self *Functions) Add(key FunctionName, val OaiFunction) {
	if self.Has(key) {
		panic(gg.Errf(
			`redundant registration of function %q of type %T`,
			key, val,
		))
	}
	self.Set(key, val)
}

func (self *Functions) Set(key FunctionName, val OaiFunction) {
	if val == nil {
		panic(gg.Errf(`unexpected nil function %q`, key))
	}
	self.OrdMap().Set(key, val)
}

func (self *Functions) OrdMap() *gg.OrdMap[FunctionName, OaiFunction] {
	return (*gg.OrdMap[FunctionName, OaiFunction])(self)
}

func (self Functions) Response(name FunctionName, arg string, verb u.Verbose) (_ string) {
	fun := self.Get(name)
	if fun == nil {
		if verb.Verb {
			log.Printf(`found no registered function %q; we consider this equivalent to empty function response`, name)
		}
		return
	}

	if verb.Verb {
		defer gg.LogTimeNow(`running function `, grepr.String(name)).LogStart().LogEnd()
	}
	return fun.OaiCall(arg)
}
