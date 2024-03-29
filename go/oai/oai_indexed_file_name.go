package oai

import (
	"_/go/u"
	"regexp"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

func ParseIndexedFileNameOpt(src string) (out IndexedFileName) {
	gg.Nop1(gg.ParseCatch(src, &out))
	return
}

func ParseIndexedFileNameValid(src string) (out IndexedFileName) {
	gg.Parse(src, &out)
	out.Validate()
	return
}

/*
Represents an indexed message file name, with decoding and encoding support.
Example names:

	0000_user_msg.md
	0000_user_request.yaml
	0001_assistant_msg.md
	0002_function_msg.yaml
*/
type IndexedFileName struct {
	Index u.FileIndex
	Role  ChatMessageRole
	Type  IndexedFileType
	Ext   string
}

// Could be implemented as `!gg.Caught(self.Validate)` but that would incur
// avoidable overhead.
func (self IndexedFileName) IsValid() bool {
	return gg.IsNotZero(self.Role) && gg.IsNotZero(self.Type)
}

// TODO: validate that role and type are known values rather than random junk.
func (self IndexedFileName) Validate() {
	if gg.IsZero(self.Role) {
		panic(gg.Errf(`invalid %T: missing role`, self))
	}
	if gg.IsZero(self.Type) {
		panic(gg.Errf(`invalid %T: missing type`, self))
	}
}

func (self IndexedFileName) ValidateRole(val ChatMessageRole) {
	if self.Role != val {
		panic(gg.Errf(
			`unexpected role in %T %q: expected role %q, found role %q`,
			self, self, val, self.Role,
		))
	}
}

func (self IndexedFileName) ValidateType(val IndexedFileType) {
	if self.Type != val {
		panic(gg.Errf(
			`unexpected type in %T %q: expected type %q, found type %q`,
			self, self, val, self.Type,
		))
	}
}

func (self IndexedFileName) IndexString() string { return self.Index.String() }
func (self IndexedFileName) IsMessage() bool     { return self.Type == IndexedFileTypeMessage }
func (self IndexedFileName) IsRequest() bool     { return self.Type == IndexedFileTypeRequest }
func (self IndexedFileName) IsEval() bool        { return self.Type == IndexedFileTypeEval }

func (self IndexedFileName) ValidString() string {
	out := self.String()
	if gg.IsZero(out) {
		panic(gg.Errf(`unable to text-encode invalid %T`, self))
	}
	return out
}

func (self IndexedFileName) String() (_ string) {
	if !self.IsValid() {
		return
	}
	return self.IndexString() + `_` + string(self.Role) + `_` + string(self.Type) + self.Ext
}

func (self IndexedFileName) MarshalText() ([]byte, error) {
	return gg.ToBytes(self.String()), nil
}

func (self *IndexedFileName) UnmarshalText(src []byte) error {
	reg := ReIndexedFileNameStrict.Get()
	mat := reg.FindSubmatch(src)

	if mat == nil {
		return gg.Errf(
			`malformed indexed file name %q; valid name must match regexp %v`,
			src, grepr.String(reg.String()),
		)
	}

	err := gg.ParseCatch(mat[1], &self.Index)
	if err != nil {
		return err
	}

	self.Role = ChatMessageRole(mat[2])
	self.Type = IndexedFileType(mat[3])
	self.Ext = string(mat[4])
	return nil
}

func (self IndexedFileName) ValidateIndex(exp int) {
	own := gg.NumConv[int](self.Index)
	if own == exp {
		return
	}

	panic(gg.Errf(
		`index mismatch in %T %q: found index %v, expected index %v`,
		self, self, own, exp,
	))
}

func ValidateIndexedFileNames(src []IndexedFileName) {
	var prev IndexedFileName

	for _, next := range src {
		next.Validate()

		if gg.IsZero(prev) {
			prev = next
			continue
		}

		// Identical index is allowed because we support having multiple files with
		// the same index (with different file extensions), automatically merging
		// them. However, such files must have the same role.
		if prev.Index == next.Index {
			if prev.Role != next.Role {
				panic(gg.Errf(
					`unexpected role mismatch between identically-indexed files %q (role %q) and %q (role %q)`,
					prev, prev.Role, next, next.Role,
				))
			}
			prev = next
			continue
		}

		// This is the normal case where 0000 is followed by 0001, etc..
		if prev.Index+1 != next.Index {
			panic(gg.Errf(
				`unexpected non-sequential file indexes: %q followed by %q`,
				prev, next,
			))
		}
		prev = next
	}
}

func IsIndexedFileNameLax(val string) bool {
	return ReIndexedFileNameLax.Get().MatchString(val)
}

var ReIndexedFileNameLax = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^\d`)
})

/*
TODO: derive digit count from `FileIndex.StringDigitCount` instead of
copy-pasting / hardcoding.
*/
var ReIndexedFileNameStrict = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^(?P<index>\d{4})_(?P<role>[a-z][a-z\d]*)_(?P<type>[a-z][a-z\d]*)(?P<ext>[.][a-z]+)?$`)
})
