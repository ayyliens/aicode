package u

import (
	"path/filepath"
	"regexp"

	"github.com/mitranim/gg"
)

func VersionedDirForkPath(path string) string {
	return ReplaceBaseName(path, VersionedDirNameInc(path).String())
}

func VersionedDirNameInc(path string) VersionedDirName {
	return gg.Max(SiblingVersionedDirNamesFrom(path)...).Inc()
}

/*
Returns a list of parsed "sibling" directory names in the parent directory of
the given directory, including its own name.

TODO cleaner code.
*/
func SiblingVersionedDirNamesFrom(path string) []VersionedDirName {
	var own VersionedDirName
	gg.Parse(filepath.Base(path), &own)

	return gg.MapCompact(
		ReadDirDirNames(filepath.Dir(path)),
		func(src string) (_ VersionedDirName) {
			var tar VersionedDirName
			gg.Parse(src, &tar)
			if tar.Base == own.Base {
				return tar
			}
			return
		},
	)
}

func VersionedDirNameFrom(path string) (out VersionedDirName) {
	gg.Parse(filepath.Base(path), &out)
	return
}

/*
Implements parsing and encoding of directory names with optional indexes:

	some_dir
	some_dir-0000
	some_dir-0001
	some_dir-0002
	...
*/
type VersionedDirName struct {
	Base    string
	Version gg.Opt[Version]
}

func (self VersionedDirName) String() (_ string) {
	if self.Version.IsNull() {
		return self.Base
	}
	return self.Base + `-` + self.VersionString()
}

func (self VersionedDirName) VersionString() string { return self.Version.String() }

func (self VersionedDirName) MarshalText() ([]byte, error) {
	return gg.ToBytes(self.String()), nil
}

func (self *VersionedDirName) UnmarshalText(src []byte) error {
	gg.PtrClear(self)

	mat := ReVersionedDirName.Get().FindSubmatch(src)
	if mat == nil {
		self.Base = string(src)
		return nil
	}

	self.Base = string(mat[1])

	err := gg.ParseCatch(mat[2], &self.Version.Val)
	if err != nil {
		return err
	}

	self.Version.Ok = true
	return nil
}

var ReVersionedDirName = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`(.*)-([\d.]+)?$`)
})

func (self VersionedDirName) GetVersion() Version { return self.Version.Val }

func (self VersionedDirName) Less(val VersionedDirName) bool {
	if self.Version.IsNull() && !val.Version.IsNull() {
		return true
	}
	return self.Version.Val.Less(val.Version.Val)
}

func (self VersionedDirName) Inc() VersionedDirName {
	if self.Version.IsNull() {
		self.Version.Set(Version{0})
		return self
	}
	// TODO maybe minor
	self.Version.Set(self.Version.Val.NextMajor())
	return self
}
