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

func (self Version) MarshalText() ([]byte, error) {
	const radix = 10
	lastIndex := gg.LastIndex(self)

	var buf []byte
	if gg.IsEmpty(self) {
		buf = []byte(FileIndex(0).String())
		return buf, nil
	}

	for ind, ver := range self {
		buf = append(buf, []byte(ver.String())...)
		if ind != lastIndex {
			buf = append(buf, `.`...)
		}
	}
	return buf, nil
}

func (self *Version) UnmarshalText(src []byte) error {
	split := gg.Split(gg.ToString(src), `.`)

	for _, ver := range split {
		var val uint16
		err := gg.ParseCatch(ver, &val)
		if err != nil {
			return err
		}
		gg.Append(self, FileIndex(val))
	}

	return nil
}
