package u_test

import (
	"_/go/u"
	"testing"

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
