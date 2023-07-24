package oai

import (
	r "reflect"

	"github.com/mitranim/gg"
)

type OaiFunctions gg.OrdMap[FunctionName, OaiFunction]

func (self OaiFunctions) Has(key FunctionName) bool {
	return self.OrdMap().Has(key)
}

func (self OaiFunctions) Get(key FunctionName) OaiFunction {
	return self.OrdMap().Get(key)
}

func (self *OaiFunctions) Add(key FunctionName, val OaiFunction) {
	if self.Has(key) {
		panic(gg.Errf(
			`redundant registration of function %q of type %T`,
			key, val,
		))
	}
	self.Set(key, val)
}

func (self *OaiFunctions) Set(key FunctionName, val OaiFunction) {
	if val == nil {
		panic(gg.Errf(`unexpected nil function %q`, key))
	}

	kind := r.TypeOf(val).Kind()
	if kind != r.Ptr {
		panic(gg.Errf(
			`unexpected non-pointer function %q of type %T; found kind: %q; function values must be pointers because we decode function call arguments into them`,
			key, val, kind,
		))
	}

	self.OrdMap().Set(key, val)
}

func (self *OaiFunctions) OrdMap() *gg.OrdMap[FunctionName, OaiFunction] {
	return (*gg.OrdMap[FunctionName, OaiFunction])(self)
}
