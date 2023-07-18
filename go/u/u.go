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
	"github.com/mitranim/gg"
	"github.com/mitranim/jsonfmt"
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

// TODO move to `gg`.
func WriteFile[A gg.Text](path string, src A) {
	gg.MkdirAll(filepath.Dir(path))
	gg.WriteFile(path, src)
}

// TODO move to `gg`. TODO better naming scheme.
func WriteFileOpt[A gg.Text](path string, src A) {
	if gg.IsNotZero(src) {
		WriteFile(path, src)
		return
	}
	if gg.FileExists(path) {
		gg.WriteFile(path, ``)
	}
}

func IsErrFileNotFound(err error) bool { return errors.Is(err, os.ErrNotExist) }

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
	ext := filepath.Ext(path)

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
	ext := filepath.Ext(path)

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

func JsonDecodeFile[A any](path string, tar *A) {
	defer gg.Detailf(`unable to decode %q as JSON`, path)
	gg.JsonDecodeFile(path, tar)
}

func JsonDecodeFileOpt[A any](path string, tar *A) {
	src := strings.TrimSpace(ReadFileOpt[string](path))
	if len(src) > 0 {
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

func YamlDecode[Src gg.Text, Tar any](src Src, tar *Tar) {
	gg.Try(yaml.Unmarshal(gg.ToBytes(src), tar))
}

func YamlDecodeOpt[Src gg.Text, Tar any](src Src, tar *Tar) {
	if len(src) > 0 {
		YamlDecode(src, tar)
	}
}

func YamlDecodeFile[A any](path string, tar *A) {
	defer gg.Detailf(`unable to decode %q as YAML`, path)
	YamlDecode(gg.ReadFile[string](path), tar)
}

func YamlDecodeFileOpt[A any](path string, tar *A) {
	src := strings.TrimSpace(ReadFileOpt[string](path))
	if len(src) > 0 {
		defer gg.Detailf(`unable to decode %q as YAML`, path)
		YamlDecode(src, tar)
	}
}

func YamlEncodeFile[A any](path string, src A) {
	defer gg.Detailf(`unable to encode %q as YAML`, path)
	gg.WriteFile(path, gg.Try1(yaml.Marshal(src)))
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
	if len(src) > 0 {
		TomlDecode(src, tar)
	}
}

func TomlDecodeFile[A any](path string, tar *A) {
	defer gg.Detailf(`unable to decode %q as TOML`, path)
	TomlDecode(gg.ReadFile[string](path), tar)
}

func TomlDecodeFileOpt[A any](path string, tar *A) {
	src := strings.TrimSpace(ReadFileOpt[string](path))
	if len(src) > 0 {
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
func ReadDirFileNames(path string) []string {
	return gg.MapCompact(gg.ReadDir(path), dirEntryToFileName)
}

func dirEntryToFileName(src fs.DirEntry) (_ string) {
	if src == nil || src.IsDir() {
		return
	}
	return src.Name()
}

type Pathed struct{ Path string }

func (self Pathed) PathJoin(path string) string {
	return filepath.Join(self.Path, path)
}

type Verbose struct{ Verb bool }

// TODO: anything built in?
func StringPadPrefix(src string, char rune, count int) string {
	var buf gg.Buf
	buf.AppendRuneN(char, count-gg.CharCount(src))
	buf.AppendString(src)
	return buf.String()
}

func FormatVerbose(src any) string {
	if src == nil {
		return ``
	}
	return strings.TrimSpace(fmt.Sprintf(`%+v`, src))
}
