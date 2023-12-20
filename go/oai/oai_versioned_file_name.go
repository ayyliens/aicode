package oai

import (
	"_/go/u"
	"regexp"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

func ParseIndexedFileNameOpt(src string) (out VersionedFileName) {
	gg.Nop1(gg.ParseCatch(src, &out))
	return
}

func ParseIndexedFileNameValid(src string) (out VersionedFileName) {
	gg.Parse(src, &out)
	out.Validate()
	return
}

/*
Represents an indexed message file name, with decoding and encoding support.
Example names:

	0_user_msg.md
	0.1_user_request.yaml
	1_assistant_msg.md
	2.0_function_msg.yaml
*/
type VersionedFileName struct {
	Version u.Version
	Role    ChatMessageRole
	Type    VersionedFileType
	Ext     string
}

// Could be implemented as `!gg.Caught(self.Validate)` but that would incur
// avoidable overhead.
func (self VersionedFileName) IsValid() bool {
	return gg.IsNotZero(self.Role) && gg.IsNotZero(self.Type)
}

// TODO: validate that role and type are known values rather than random junk.
func (self VersionedFileName) Validate() {
	if gg.IsZero(self.Role) {
		panic(gg.Errf(`invalid %T: missing role`, self))
	}
	if gg.IsZero(self.Type) {
		panic(gg.Errf(`invalid %T: missing type`, self))
	}
}

func (self VersionedFileName) ValidateRole(val ChatMessageRole) {
	if self.Role != val {
		panic(gg.Errf(
			`unexpected role in %T %q: expected role %q, found role %q`,
			self, self, val, self.Role,
		))
	}
}

func (self VersionedFileName) ValidateType(val VersionedFileType) {
	if self.Type != val {
		panic(gg.Errf(
			`unexpected type in %T %q: expected type %q, found type %q`,
			self, self, val, self.Type,
		))
	}
}

func (self VersionedFileName) VersionString() string { return self.Version.String() }
func (self VersionedFileName) IsMessage() bool       { return self.Type == VersionedFileTypeMessage }
func (self VersionedFileName) IsRequest() bool       { return self.Type == VersionedFileTypeRequest }
func (self VersionedFileName) IsEval() bool          { return self.Type == VersionedFileTypeEval }

func (self VersionedFileName) GetVersion() u.Version    { return self.Version }
func (self VersionedFileName) GetRole() ChatMessageRole { return self.Role }
func (self VersionedFileName) GetIndexRole() string     { return gg.Str(self.Version, self.Role) }

func (self VersionedFileName) ValidString() string {
	out := self.String()
	if gg.IsZero(out) {
		panic(gg.Errf(`unable to text-encode invalid %T`, self))
	}
	return out
}

/*
TODO sync symbol for version ending as in regex ReIndexedFileNameStrict
*/
func (self VersionedFileName) Name() (_ string) {
	if !self.IsValid() {
		return
	}
	return self.VersionString() + `-` + string(self.Role) + `-` + string(self.Type)
}

func (self VersionedFileName) String() (_ string) {
	if !self.IsValid() {
		return
	}
	return self.Name() + self.Ext
}

func (self VersionedFileName) MarshalText() ([]byte, error) {
	return gg.ToBytes(self.String()), nil
}

func (self *VersionedFileName) UnmarshalText(src []byte) error {
	reg := ReIndexedFileNameStrict.Get()
	mat := reg.FindSubmatch(src)

	if mat == nil {
		return gg.Errf(
			`malformed indexed file name %q; valid name must match regexp %v`,
			src, grepr.String(reg.String()),
		)
	}

	err := gg.ParseCatch(mat[1], &self.Version)
	if err != nil {
		return nil
	}

	self.Role = ChatMessageRole(mat[2])
	self.Type = VersionedFileType(mat[3])
	self.Ext = string(mat[4])
	return nil
}

func (self VersionedFileName) Less(val VersionedFileName) bool {
	if self.Version.Equal(val.Version) {
		return self.Role.Index() < val.Role.Index()
	}

	return self.Version.Less(val.Version)
}

func ValidateIndexedFileNames(src []VersionedFileName) {
	var prev VersionedFileName

	for _, next := range src {
		next.Validate()

		if gg.IsZero(prev) {
			prev = next
			continue
		}

		// Identical index is allowed because we support having multiple files with
		// the same index (with different file extensions or roles), automatically merging
		// them.
		if prev.Version.Equal(next.Version) {
			prev = next
			continue
		}

		// This is the normal case where 0000 is followed by 0001, etc..
		if !prev.Version.IsNext(next.Version) {
			panic(gg.Errf(
				`unexpected non-sequential file indexes: %q followed by %q`,
				prev, next,
			))
		}
		prev = next
	}
}

func IsVersionedFileNameLax(val string) bool {
	return ReIndexedFileNameLax.Get().MatchString(val)
}

var ReIndexedFileNameLax = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^\d`)
})

var ReIndexedFileNameStrict = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^(?P<version>[\d.]+)-(?P<role>[a-z][a-z\d]*)-(?P<type>[a-z][a-z\d]*)(?P<ext>[.][a-z]+)?$`)
})
