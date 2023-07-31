package oai

import (
	"github.com/mitranim/gg"
)

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
