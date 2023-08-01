package u_test

import (
	"_/go/u"
	"strings"
	"testing"

	"github.com/mitranim/gg/gtest"
)

func Test_YamlJsonString(t *testing.T) {
	defer gtest.Catch(t)

	src := strings.TrimSpace(`
one:
  - two
  - three
`)

	var tar u.YamlJsonString
	u.YamlDecode(src, &tar)

	gtest.Eq(tar, `{"one":["two","three"]}`)
	gtest.Eq(strings.TrimSpace(u.YamlEncode[string](tar)), src)
}
