package u_test

import (
	"_/go/u"
	"math"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func Test_StringPadPrefix(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(u.StringPadPrefix(`1`, '0', 4), `0001`)
	gtest.Eq(u.StringPadPrefix(`12`, '0', 4), `0012`)
	gtest.Eq(u.StringPadPrefix(`123`, '0', 4), `0123`)
	gtest.Eq(u.StringPadPrefix(`1234`, '0', 4), `1234`)
	gtest.Eq(u.StringPadPrefix(`12345`, '0', 4), `12345`)
}

func Test_FilepathDir(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(filepath.Dir(`/one/two`), `/one`)
	gtest.Eq(u.FilepathDir(`/one/two`), `/one`)

	gtest.Eq(filepath.Dir(`/one`), `/`)
	gtest.Eq(u.FilepathDir(`/one`), `/`)

	gtest.Eq(filepath.Dir(`/`), `/`)
	gtest.Eq(u.FilepathDir(`/`), ``) // Difference here.
}

func Test_PkgRoot(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(
		gg.Try1(filepath.Rel(u.PkgRoot(), gg.Cwd())),
		`go/u`,
	)
}

func Test_PkgRelPath(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(strings.HasPrefix(
		u.PkgRelPath(`some_dir/some_file`),
		u.PkgRoot(),
	))
}

func Test_FileExt(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(u.FileExt(``), ``)

	gtest.Eq(u.FileExt(`one`), ``)
	gtest.Eq(u.FileExt(`.one`), ``)
	gtest.Eq(u.FileExt(`one.`), `.`)
	gtest.Eq(u.FileExt(`one.two`), `.two`)

	gtest.Eq(u.FileExt(`one/two`), ``)
	gtest.Eq(u.FileExt(`one/two.`), `.`)
	gtest.Eq(u.FileExt(`one/two.three`), `.three`)
	gtest.Eq(u.FileExt(`one.two/three`), ``)
	gtest.Eq(u.FileExt(`one.two/.three`), ``)
	gtest.Eq(u.FileExt(`one.two/three.`), `.`)
	gtest.Eq(u.FileExt(`one.two/three.four`), `.four`)
}

func Test_filepath_Join_appending_absolute_path(t *testing.T) {
	defer gtest.Catch(t)

	const baseDirName = `/tmp/os_temp_dir`
	const tempDirName = `tool_name`
	const pathSuffix = `/one/two/three`

	tar := filepath.Join(baseDirName, tempDirName, pathSuffix)

	gtest.Eq(
		tar,
		`/tmp/os_temp_dir/tool_name/one/two/three`,
	)
}

func Test_IntStringDigitCount(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`invalid radix 0 for numeric encoding`,
		func() { u.IntStringDigitCount[byte](234, 0) },
	)

	gtest.PanicStr(
		`invalid radix 1 for numeric encoding`,
		func() { u.IntStringDigitCount[byte](234, 1) },
	)

	testIntStringDigitCount[byte](0, 10, 1)
	testIntStringDigitCount[byte](1, 10, 1)
	testIntStringDigitCount[byte](2, 10, 1)
	testIntStringDigitCount[byte](9, 10, 1)
	testIntStringDigitCount[byte](10, 10, 2)
	testIntStringDigitCount[byte](99, 10, 2)
	testIntStringDigitCount[byte](100, 10, 3)
	testIntStringDigitCount[byte](255, 10, 3)

	testIntStringDigitCount[int8](0, 10, 1)
	testIntStringDigitCount[int8](1, 10, 1)
	testIntStringDigitCount[int8](2, 10, 1)
	testIntStringDigitCount[int8](9, 10, 1)
	testIntStringDigitCount[int8](10, 10, 2)
	testIntStringDigitCount[int8](99, 10, 2)
	testIntStringDigitCount[int8](100, 10, 3)
	testIntStringDigitCount[int8](127, 10, 3)
	testIntStringDigitCount[int8](-1, 10, 2)
	testIntStringDigitCount[int8](-9, 10, 2)
	testIntStringDigitCount[int8](-10, 10, 3)
	testIntStringDigitCount[int8](-99, 10, 3)
	testIntStringDigitCount[int8](-100, 10, 4)
	testIntStringDigitCount[int8](-128, 10, 4)

	testIntStringDigitCount[uint64](math.MaxUint64, 2, 64)
	testIntStringDigitCount[uint64](math.MaxUint64, 10, 20)
	testIntStringDigitCount[uint64](math.MaxUint64, 16, 16)
	testIntStringDigitCount[uint64](math.MaxUint64, 32, 13)

	testIntStringDigitCount[int64](math.MaxInt64, 2, 63)
	testIntStringDigitCount[int64](math.MaxInt64, 10, 19)
	testIntStringDigitCount[int64](math.MaxInt64, 16, 16)
	testIntStringDigitCount[int64](math.MaxInt64, 32, 13)
	testIntStringDigitCount[int64](math.MinInt64, 2, 65)
	testIntStringDigitCount[int64](math.MinInt64, 10, 20)
	testIntStringDigitCount[int64](math.MinInt64, 16, 17)
	testIntStringDigitCount[int64](math.MinInt64, 32, 14)
}

func testIntStringDigitCount[A gg.Int](src A, rad, exp byte) {
	str := u.FormatInt(src, rad)

	msg := gg.JoinLinesOpt(
		gg.Str(`type: `, gg.Type[A]()),
		gg.Str(`number: `, src),
		gg.Str(`radix: `, rad),
		gg.Str(`encoded: `, str),
	)

	gtest.Eq(u.IntStringDigitCount(src, rad), exp, msg)
	gtest.Eq(len(str), int(exp), msg)
}
