package u

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/mitranim/gg"
	"github.com/mitranim/jsonfmt"
	mofo "golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

type Ctx = context.Context

func ParseJsonLines[Tar any, Src gg.Text](src Src) []Tar {
	return DecodeJsonLines[Tar](gg.NewReadCloser(src))
}

func DecodeJsonLines[Tar any, Src io.Reader](src Src) (out []Tar) {
	DecodeJsonLinesInto(src, &out)
	return
}

func DecodeJsonLinesInto[Slice ~[]Elem, Elem any, Src io.Reader](src Src, out *Slice) {
	dec := json.NewDecoder(src)
	for dec.More() {
		gg.Try(dec.Decode(gg.AppendPtrZero(out)))
	}
	return
}

func WriteJsonLines[A any](path string, src []A) {
	tar := gg.Try1(os.Create(path))
	defer gg.Close(tar)

	enc := json.NewEncoder(tar)
	for _, val := range src {
		gg.Try(enc.Encode(val))
	}
}

func FileRead[A gg.Text](src *os.File) A {
	gg.Try1(src.Seek(0, 0))
	return gg.ToText[A](gg.Try1(io.ReadAll(src)))
}

func FileRewrite[A gg.Text](tar *os.File, src A) {
	gg.Try1(tar.Seek(0, 0))
	gg.Try1(tar.Write(gg.ToBytes(src)))
}

// "Rec" stands for "recursive" (with mkdir). TODO better naming scheme.
func TouchedDirRec(path string) bool {
	if gg.DirExists(path) {
		return false
	}
	gg.MkdirAll(path)
	return true
}

// "Rec" stands for "recursive" (with mkdir). TODO better naming scheme.
func TouchedFileRec(path string) bool {
	dir := TouchedDirRec(filepath.Dir(path))
	file := TouchedFile(path)
	return dir || file
}

func TouchedFile(path string) bool {
	if gg.FileExists(path) {
		return false
	}
	gg.Try1(os.OpenFile(path, os.O_CREATE, os.ModePerm)).Close()
	return true
}

func ReadFileOpt[A gg.Text](path string) A {
	defer gg.SkipOnly(IsErrFileNotFound)
	return gg.ReadFile[A](path)
}

/*
Stands for "copy directory recursive". Caution: parameter order is inconsistent
with built-in `copy` and `io.Copy`.

TODO: validate that the source and destination directories are not ancestor or
descendant of each other.
*/
func CopyDirRec(src, out string) {
	var exists bool

	for _, entry := range ReadDirOpt(src) {
		exists = exists || TouchedDirRec(out)

		srcPath := filepath.Join(src, entry.Name())
		outPath := filepath.Join(out, entry.Name())

		if entry.IsDir() {
			CopyDirRec(srcPath, outPath)
		} else {
			CopyFile(srcPath, outPath)
		}
	}
}

// Caution: parameter order is inconsistent with built-in `copy` and `io.Copy`.
func CopyFile(srcPath, outPath string) {
	src := gg.Try1(os.Open(srcPath))
	defer gg.Close(src)

	out := gg.Try1(os.Create(outPath))
	defer gg.Close(out)

	gg.Try1(io.Copy(out, src))
}

func IsErrFileNotFound(err error) bool { return errors.Is(err, os.ErrNotExist) }

func RemoveFileOrDir(path string) { gg.Try(os.Remove(path)) }

func RemoveFileOrDirOpt(path string) { gg.Nop1(os.Remove(path)) }

func LogErr(err error) {
	if err == nil {
		return
	}
	log.Printf(`error: %+v`, err)
}

func Wait(ctx Ctx, dur time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(dur):
	}
}

func PolyDecodeFileOpt[A any](path string, tar *A) {
	ext := FileExt(path)

	switch ext {
	case `.json`:
		JsonDecodeFileOpt(path, tar)
	case `.yaml`:
		YamlDecodeFileOpt(path, tar)
	case `.toml`:
		TomlDecodeFileOpt(path, tar)
	default:
		panic(gg.Errf(`unable to polymorphically decode %q: unrecognized extension %q`, path, ext))
	}
}

func PolyEncodeFileOpt[A any](path string, src A) {
	ext := FileExt(path)

	switch ext {
	case `.json`:
		JsonEncodeFileOpt(path, src)
	case `.yaml`:
		YamlEncodeFileOpt(path, src)
	case `.toml`:
		TomlEncodeFileOpt(path, src)
	default:
		panic(gg.Errf(`unable to polymorphically encode %q: unrecognized extension %q`, path, ext))
	}
}

func JsonPretty[A gg.Text](src A) A {
	return gg.ToText[A](jsonfmt.FormatString(jsonfmt.Default, src))
}

func JsonEncodePretty[Tar gg.Text, Src any](src Src) Tar {
	return JsonPretty(gg.JsonEncode[Tar](src))
}

/*
Same as `json.Unmarshal` but with panics and support for arbitrary source text
types. Same as `gg.JsonDecode` but takes `any` instead of explicit pointer.
TODO move to `gg`.
*/
func JsonDecodeAny[A gg.Text](src A, out any) {
	gg.Try(json.Unmarshal(gg.ToBytes(src), out))
}

func JsonDecodeFile[A any](path string, tar *A) {
	defer gg.Detailf(`unable to decode %q as JSON`, path)
	gg.JsonDecodeFile(path, tar)
}

func JsonDecodeFileOpt[A any](path string, tar *A) {
	src := strings.TrimSpace(ReadFileOpt[string](path))
	if gg.IsNotZero(src) {
		defer gg.Detailf(`unable to decode %q as JSON`, path)
		gg.JsonDecode(src, tar)
	}
}

func JsonEncodeFile[A any](path string, src A) {
	defer gg.Detailf(`unable to encode %q as JSON`, path)
	WriteFile(path, JsonEncodePretty[string](src))
}

func JsonEncodeFileOpt[A any](path string, src A) {
	if gg.IsNotZero(src) {
		JsonEncodeFile(path, src)
		return
	}
	if gg.FileExists(path) {
		gg.WriteFile(path, ``)
	}
}

// Difference from `yaml.Marshal`: indent two spaces.
func YamlEncode[Tar gg.Text, Src any](src Src) Tar {
	var buf gg.Buf
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	gg.Try(enc.Encode(src))
	return gg.ToText[Tar](buf)
}

func YamlDecode[Src gg.Text, Tar any](src Src, tar *Tar) {
	gg.Try(yaml.Unmarshal(gg.ToBytes(src), tar))
}

func YamlDecodeOpt[Src gg.Text, Tar any](src Src, tar *Tar) {
	if gg.IsNotZero(src) {
		YamlDecode(src, tar)
	}
}

func YamlDecodeFile[A any](path string, tar *A) {
	defer gg.Detailf(`unable to decode %q as YAML`, path)
	YamlDecode(gg.ReadFile[string](path), tar)
}

func YamlDecodeFileOpt[A any](path string, tar *A) {
	src := strings.TrimSpace(ReadFileOpt[string](path))
	if gg.IsNotZero(src) {
		defer gg.Detailf(`unable to decode %q as YAML`, path)
		YamlDecode(src, tar)
	}
}

func YamlEncodeFile[A any](path string, src A) {
	defer gg.Detailf(`unable to encode %q as YAML`, path)
	gg.WriteFile(path, YamlEncode[string](src))
}

func YamlEncodeFileOpt[A any](path string, src A) {
	if gg.IsNotZero(src) {
		YamlEncodeFile(path, src)
		return
	}
	if gg.FileExists(path) {
		gg.WriteFile(path, ``)
	}
}

func TomlDecode[Src gg.Text, Tar any](src Src, tar *Tar) {
	gg.Try(toml.Unmarshal(gg.ToBytes(src), tar))
}

func TomlDecodeOpt[Src gg.Text, Tar any](src Src, tar *Tar) {
	if gg.IsNotZero(src) {
		TomlDecode(src, tar)
	}
}

func TomlDecodeFile[A any](path string, tar *A) {
	defer gg.Detailf(`unable to decode %q as TOML`, path)
	TomlDecode(gg.ReadFile[string](path), tar)
}

func TomlDecodeFileOpt[A any](path string, tar *A) {
	src := strings.TrimSpace(ReadFileOpt[string](path))
	if gg.IsNotZero(src) {
		defer gg.Detailf(`unable to decode %q as TOML`, path)
		TomlDecode(src, tar)
	}
}

// Note: the TOML package does not provide a "marshal" function.
func TomlEncode[Out gg.Text, Src any](src Src) Out {
	var buf gg.Buf
	gg.Try(toml.NewEncoder(&buf).Encode(src))
	return gg.ToText[Out](buf)
}

func TomlEncodeFile[A any](path string, src A) {
	defer gg.Detailf(`unable to encode %q as TOML`, path)
	gg.WriteFile(path, TomlEncode[string](src))
}

func TomlEncodeFileOpt[A any](path string, src A) {
	if gg.IsNotZero(src) {
		TomlEncodeFile(path, src)
		return
	}
	if gg.FileExists(path) {
		gg.WriteFile(path, ``)
	}
}

// TODO move to `gg`.
func ReadDirOpt(path string) []fs.DirEntry {
	defer gg.SkipOnly(IsErrFileNotFound)
	return gg.ReadDir(path)
}

// TODO move to `gg`.
func ReadDirFileNames(path string) []string {
	return gg.MapCompact(ReadDirOpt(path), dirEntryToFileName)
}

func dirEntryToFileName(src fs.DirEntry) (_ string) {
	if src == nil || src.IsDir() {
		return
	}
	return src.Name()
}

// TODO better name. TODO move to `gg`.
func ReadDirDirNames(path string) []string {
	return gg.MapCompact(ReadDirOpt(path), dirEntryToDirName)
}

func dirEntryToDirName(src fs.DirEntry) (_ string) {
	if src == nil || !src.IsDir() {
		return
	}
	return src.Name()
}

type Pathed struct{ Path string }

func (self Pathed) PathJoin(path string) string {
	return filepath.Join(self.Path, path)
}

type Verbose struct{ Verb bool }

type Inited struct{ Init bool }

// TODO: anything built in?
func StringPadPrefix(src string, char rune, count int) string {
	var buf gg.Buf
	buf.AppendRuneN(char, count-gg.CharCount(src))
	buf.AppendString(src)
	return buf.String()
}

func NumToPaddedString[A gg.Num](src A) string {
	return StringPadPrefix(gg.String(src), '0', 4)
}

func FormatVerbose(src any) string {
	if src == nil {
		return ``
	}
	return strings.TrimSpace(fmt.Sprintf(`%+v`, src))
}

/*
Difference from `filepath.Dir`: returns zero value when directory has no
parent.
*/
func FilepathDir(src string) (_ string) {
	src = filepath.Clean(src)
	out := filepath.Dir(src)
	if out == src {
		return
	}
	return out
}

func DirProcureAnc[A any](dir string, fun func(string) A) (_ string, _ A) {
	if fun == nil {
		return
	}

	dir = filepath.Clean(dir)

	for gg.IsNotZero(dir) {
		val := fun(dir)
		if gg.IsNotZero(val) {
			return dir, val
		}
		dir = FilepathDir(dir)
	}

	return
}

func DirFindAnc[A any](dir string, fun func(string) A) string {
	dir, _ = DirProcureAnc(dir, fun)
	return dir
}

func PkgRoot() string {
	out := DirFindAnc(gg.Cwd(), IsPkgRoot)
	if gg.IsZero(out) {
		panic(gg.Errv(`unable to find path of root package`))
	}
	return out
}

func IsPkgRoot(dir string) bool {
	file := ReadGomodOpt(filepath.Join(dir, `go.mod`))
	return file != nil && file.Module != nil && file.Module.Mod.Path == `_`
}

func ReadGomodOpt(path string) *mofo.File {
	return ParseGomodOpt(path, ReadFileOpt[[]byte](path))
}

func ParseGomodOpt(path string, body []byte) *mofo.File {
	if gg.IsTextEmpty(body) {
		return nil
	}
	return gg.Try1(mofo.Parse(path, body, nil))
}

var PkgRootOnce = gg.NewLazy(PkgRoot)

/*
Joins given path with repo root path. Useful for tests because `go test` changes
the current working directory when running tests in sub-folders.
*/
func PkgRelPath(path string) string {
	return filepath.Join(PkgRootOnce.Get(), path)
}

// TODO better name.
func JoinLines2Opt(src ...string) string { return gg.JoinOpt(src, "\n\n") }

func LoadEnvFiles() {
	for _, base := range gg.Reversed(strings.Split(os.Getenv(`CONF`), `,`)) {
		gg.Try(godotenv.Load(filepath.Join(base, `.env.properties`)))
	}
}

// TODO move to `gg`.
func IsTextBlank[A gg.Text](src A) bool {
	return gg.IsTextEmpty(strings.TrimSpace(gg.ToString(src)))
}

/*
Must be deferred. Same as `gg.Fail` but the given function is nullary.
TODO better name and move to `gg`.
*/
func Fail0(fun func()) {
	err := gg.AnyErrTracedAt(recover(), 1)
	if err != nil && fun != nil {
		fun()
	}
	gg.Try(err)
}

/*
Workaround for the issue where `filepath.Ext` incorrectly reports non-empty
extensions for file names such as `.blah` where `.blah` is the file name,
not the extension.

TODO: make this work for both Unix and Windows paths.
*/
func FileExt(src string) string {
	name := filepath.Base(src)
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	if base == `` {
		return ``
	}
	return ext
}

func ReplaceBaseName(src, name string) string {
	return filepath.Join(filepath.Dir(src), name)
}

func IsErrContextCancel(err error) bool {
	return errors.Is(err, context.Canceled)
}
