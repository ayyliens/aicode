package u

import (
	"strconv"

	"github.com/mitranim/gg"
)

/*
Represents a fixed-length numeric index in a file name, such as `0000`, `0001`,
and so on. Implements decoding and encoding.
*/
type FileIndex uint16

/*
Digit count must be fixed to ensure that ordering file names by the common
string-sorting algorithm is identical to ordering file names by the indexes as
integers (assuming no duplicate indexes). We COULD internally order by parsed
indexes, but we also want to ensure that files are ordered the same way in all
FS browsers, including the OS built-ins and file lists in code editors, which
requires a fixed count.
*/
func (FileIndex) StringDigitCount() int { return 4 }

func (self FileIndex) String() string {
	return gg.ToString(gg.Try1(self.MarshalText()))
}

func (self FileIndex) MarshalText() ([]byte, error) {
	const radix = 10
	expCount := self.StringDigitCount()
	ownCount := int(IntStringDigitCount(self, radix))
	missing := expCount - ownCount

	if missing < 0 {
		return nil, gg.Errf(
			`%T %v overflows allowed digit count %v in radix %v`,
			self, int(self), expCount, radix,
		)
	}

	buf := make([]byte, expCount)
	for ind := range buf[:missing] {
		buf[ind] = '0'
	}

	// Ideally, we would verify that the number of appended bytes is equal to
	// `missing`.
	strconv.AppendUint(buf[missing:missing], uint64(self), radix)

	return buf, nil
}

func (self *FileIndex) UnmarshalText(src []byte) error {
	expCount := self.StringDigitCount()
	srcCount := len(src)

	if expCount != srcCount {
		return gg.Errf(
			`unable to decode %q as %T: length mismatch: expected %v digits, found %v digits`,
			src, self, expCount, srcCount,
		)
	}

	return gg.ParseCatch(src, (*uint16)(self))
}
