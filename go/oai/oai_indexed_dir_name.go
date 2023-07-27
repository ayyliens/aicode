package oai

import (
	"_/go/u"
	"regexp"

	"github.com/mitranim/gg"
)

/*
Note: the amount of digits must be fixed. See the comment on
`ReMessageFileNameStrict` for an explanation.
*/
var ReIndexedDirName = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`(.*)_(\d{4})?$`)
})

type IndexedDirName struct {
	Base  string
	Index gg.Opt[uint]
}

func (self IndexedDirName) String() (_ string) {
	if self.Index.IsNull() {
		return self.Base
	}
	return self.Base + `_` + self.IndexString()
}

func (self IndexedDirName) IndexString() string {
	return u.NumToPaddedString(self.Index.Val)
}

func (self *IndexedDirName) Parse(src string) error {
	gg.PtrClear(self)

	mat := ReIndexedDirName.Get().FindStringSubmatch(src)
	if mat == nil {
		self.Base = src
		return nil
	}

	self.Base = mat[1]
	gg.Parse(mat[2], &self.Index.Val)
	self.Index.Ok = true
	return nil
}

// func (self IndexedDirName) ValidateIndex(ind int) {
// 	ValidateIndex[IndexedDirName](gg.NumConv[int](self.Index.Val), ind)
// }

func (self IndexedDirName) GetIndex() uint { return self.Index.Val }

func (self IndexedDirName) Less(val IndexedDirName) bool {
	if self.Index.IsNull() && !val.Index.IsNull() {
		return true
	}
	return self.Index.Val < val.Index.Val
}

func (self IndexedDirName) Inc() IndexedDirName {
	if self.Index.IsNull() {
		self.Index.Set(0)
		return self
	}
	self.Index.Set(gg.Inc(self.Index.Val))
	return self
}
