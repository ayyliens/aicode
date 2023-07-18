package oai

import (
	"_/go/u"
	"regexp"
	"strings"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

func IsMsgFileNameLax(val string) bool {
	return ReMsgFileNameLax.Get().MatchString(val)
}

var ReMsgFileNameLax = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^msg_\d+_`)
})

var ReMsgFileNameStrict = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^msg_(\d+)_([a-z][a-z\d]*)([.][a-z]+)?$`)
})

type MsgFileName struct {
	Index int
	Role  ChatMessageRole
	Ext   string
}

func (self MsgFileName) IsValid() bool { return gg.IsNotZero(self.Role) }

func (self MsgFileName) String() (_ string) {
	if !self.IsValid() {
		return
	}
	return `msg_` + self.IndexString() + `_` + string(self.Role) + self.Ext
}

func (self MsgFileName) IndexString() (_ string) {
	return u.StringPadPrefix(gg.String(self.Index), '0', 4)
}

func (self *MsgFileName) Parse(src string) (err error) {
	defer gg.Rec(&err)

	reg := ReMsgFileNameStrict.Get()
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

func (self MsgFileName) ValidateIndex(ind int) {
	if ind == self.Index {
		return
	}

	panic(gg.Errf(
		`index mismatch in %T: expected %v, found %v`,
		self, ind, self.Index,
	))
}

// TODO consider validating the message.
func (self MsgFileName) ChatCompletionMessage(path string) (out ChatCompletionMessage) {
	defer gg.Detailf(`unable to decode msg %q`, path)

	out.Role = self.Role

	switch self.Ext {
	case ``, `.txt`, `.md`:
		out.Content = strings.TrimSpace(gg.ReadFile[string](path))
		return

	case `.json`:
		u.JsonDecodeFileOpt(path, &out)
		MsgValidateRoleMatch(path, out.Role, self.Role)
		return

	case `.yaml`:
		u.YamlDecodeFileOpt(path, &out)
		MsgValidateRoleMatch(path, out.Role, self.Role)
		return

	case `.toml`:
		u.TomlDecodeFileOpt(path, &out)
		MsgValidateRoleMatch(path, out.Role, self.Role)
		return

	default:
		panic(gg.Errf(`unrecognized msg file extension %q`, self.Ext))
	}
}

func MsgValidateRoleMatch(path string, act, exp ChatMessageRole) {
	if gg.NotEqNotZero(act, exp) {
		panic(gg.Errf(
			`unexpected role mismatch in msg %q: expected be %q or empty, found %q`,
			path, exp, act,
		))
	}
}
