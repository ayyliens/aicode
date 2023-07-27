package u_test

import (
	"_/go/u"
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
