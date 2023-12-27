package u

import (
	"github.com/mitranim/gg"
)

/*
Represents a version in a file name, such as `0.0.0.0`, `1` == `1.0.0.0`, `1.0.0.1`
and so on. Implements decoding and encoding.
*/
type Version []FileIndex

func (self Version) String() string {
	return gg.ToString(gg.Try1(self.MarshalText()))
}

func (self Version) IsNext(val Version) bool {
	maxLen := gg.MaxPrim(len(self), len(val))
	for i := 0; i <= maxLen; i++ {
		a := gg.Get(self, i)
		b := gg.Get(val, i)
		if a+1 == b {
			return true
		} else if a != b {
			return false
		}
	}

	return false
}

func (self Version) Equal(val Version) bool {
	maxLen := gg.MaxPrim(len(self), len(val))
	for i := 0; i <= maxLen; i++ {
		a := gg.Get(self, i)
		b := gg.Get(val, i)
		if a != b {
			return false
		}
	}
	return true
}

func (self Version) Less(val Version) bool {
	maxLen := gg.MaxPrim(len(self), len(val))
	for i := 0; i <= maxLen; i++ {
		a := gg.Get(self, i)
		b := gg.Get(val, i)
		if a < b {
			return true
		} else if a > b {
			return false
		}
	}
	return false
}

func (self Version) NextMajor() Version {
	return Version{gg.Head(self) + 1}
}

func (self Version) NextMinor() Version {
	return append(gg.Init(self), gg.Last(self)+1)
}

func (self Version) PrevMajor(depth uint16) Version {
	return Version{FileIndex(uint16(gg.Head(self)) - depth)}
}

func (self Version) AddMinor() Version {
	return append(self, 1)
}

func (self Version) MarshalText() (out []byte, err error) {
	return self.AppendTo(out), nil
}

func (self *Version) UnmarshalText(src []byte) error {
	return self.Parse(gg.ToString(src))
}

func (self Version) AppendTo(buf []byte) []byte {
	val := gg.Join(gg.Map(self, FileIndex.String), `.`)
	return append(buf, gg.ToBytes(val)...)
}

func (self *Version) Parse(src string) (err error) {
	defer gg.Rec(&err)
	gg.Append(self, gg.Map(gg.Split(src, `.`), gg.Zero[FileIndex]().TryParse)...)
	return
}
