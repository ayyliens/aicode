package oai

import (
	"_/go/u"
	"regexp"
	"strings"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

func IsIndexedMessageFileNameLax(val string) bool {
	return ReIndexedMessageFileNameLax.Get().MatchString(val)
}

var ReIndexedMessageFileNameLax = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^\d`)
})

/*
Note: the amount of digits that denote the index should be fixed, to ensure that
ordering file names by the common string-sorting algorithm is identical to
ordering file names by the indexes as integers (assuming no duplicate indexes).
We COULD internally order by parsed indexes, but we also want to ensure that
files are ordered the same way in all FS browsers, including the OS built-ins
and file lists in code editors, which requires a fixed digit count.
*/
var ReIndexedMessageFileNameStrict = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^(\d{4})_msg_([a-z][a-z\d]*)([.][a-z]+)?$`)
})

func IndexedMessageFileNameOpt(src string) (out IndexedMessageFileName) {
	gg.Nop1(out.Parse(src))
	return
}

/*
Represents an indexed message file name, with decoding and encoding support.
Example names:

	0000_msg_user.md
	0001_msg_assistant.md
	0002_msg_function.yaml
*/
type IndexedMessageFileName struct {
	Index uint
	Role  ChatMessageRole
	Ext   string
}

func (self IndexedMessageFileName) IsValid() bool { return gg.IsNotZero(self.Role) }

func (self IndexedMessageFileName) String() (_ string) {
	if !self.IsValid() {
		return
	}
	return self.IndexString() + `_msg_` + string(self.Role) + self.Ext
}

func (self IndexedMessageFileName) ValidString() string {
	out := self.String()
	if gg.IsZero(out) {
		panic(gg.Errf(`unable to string-encode invalid %T`, self))
	}
	return out
}

func (self IndexedMessageFileName) IndexString() string {
	return u.NumToPaddedString(self.Index)
}

func (self *IndexedMessageFileName) Parse(src string) (err error) {
	defer gg.Rec(&err)

	reg := ReIndexedMessageFileNameStrict.Get()
	mat := reg.FindStringSubmatch(src)

	if mat == nil {
		panic(gg.Errf(
			`malformed msg file name %q; valid name must match regexp %v`,
			src, grepr.String(reg.String()),
		))
	}

	gg.Parse(mat[1], &self.Index)
	self.Role = ChatMessageRole(mat[2])
	self.Ext = mat[3]
	return
}

func (self IndexedMessageFileName) ValidateIndex(ind int) {
	ValidateIndex[IndexedMessageFileName](gg.NumConv[int](self.Index), ind)
}

// TODO consider validating the message.
func (self IndexedMessageFileName) ChatCompletionMessageExt(path string) (out ChatCompletionMessageExt) {
	defer gg.Detailf(`unable to decode msg %q`, path)

	out.Role = self.Role
	out.FileName = self

	switch self.Ext {
	case ``, `.txt`, `.md`:
		out.Content = strings.TrimSpace(gg.ReadFile[string](path))
		return

	case `.json`:
		u.JsonDecodeFileOpt(path, &out)
		MessageValidateRoleMatch(path, out.Role, self.Role)
		return

	case `.yaml`:
		u.YamlDecodeFileOpt(path, &out)
		MessageValidateRoleMatch(path, out.Role, self.Role)
		return

	case `.toml`:
		u.TomlDecodeFileOpt(path, &out)
		MessageValidateRoleMatch(path, out.Role, self.Role)
		return

	default:
		panic(gg.Errf(`unrecognized msg file extension %q`, self.Ext))
	}
}

func MessageValidateRoleMatch(path string, act, exp ChatMessageRole) {
	if gg.NotEqNotZero(act, exp) {
		panic(gg.Errf(
			`unexpected role mismatch in msg %q: expected be %q or empty, found %q`,
			path, exp, act,
		))
	}
}

func ValidateIndex[Tar any, Num gg.Int](act, exp Num) {
	if act == exp {
		return
	}

	panic(gg.Errf(
		`index mismatch in %v: found %v, expected %v`,
		gg.Type[Tar](), act, exp,
	))
}
