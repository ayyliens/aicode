package u

import (
	"path/filepath"
	"regexp"

	"github.com/mitranim/gg"
)

func IndexedDirForkPath(path string) string {
	return ReplaceBaseName(path, IndexedDirNameInc(path).String())
}

func IndexedDirNameInc(path string) IndexedDirName {
	return gg.Max(SiblingIndexedDirNamesFrom(path)...).Inc()
}

/*
Returns a list of parsed "sibling" directory names in the parent directory of
the given directory, including its own name.

TODO cleaner code.
*/
func SiblingIndexedDirNamesFrom(path string) []IndexedDirName {
	var own IndexedDirName
	gg.Try(own.Parse(filepath.Base(path)))

	return gg.MapCompact(
		ReadDirDirNames(filepath.Dir(path)),
		func(src string) (_ IndexedDirName) {
			var tar IndexedDirName
			gg.Try(tar.Parse(src))
			if tar.Base == own.Base {
				return tar
			}
			return
		},
	)
}

func IndexedDirNameFrom(path string) (out IndexedDirName) {
	gg.Try(out.Parse(filepath.Base(path)))
	return
}

/*
Implements parsing and encoding of directory names with optional indexes:

	some_dir
	some_dir_0000
	some_dir_0001
	some_dir_0002
	...
*/
type IndexedDirName struct {
	Base  string
	Index gg.Opt[FileIndex]
}

func (self IndexedDirName) String() (_ string) {
	if self.Index.IsNull() {
		return self.Base
	}
	return self.Base + `_` + self.IndexString()
}

func (self IndexedDirName) IndexString() string { return self.Index.String() }

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

/*
Note: the amount of digits must be fixed. See the comment on
`FileIndex.StringDigitCount` for an explanation.

TODO: derive digit count from `FileIndex.StringDigitCount` instead of
copy-pasting / hardcoding.
*/
var ReIndexedDirName = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`(.*)_(\d{4})?$`)
})

func (self IndexedDirName) GetIndex() FileIndex { return self.Index.Val }

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
