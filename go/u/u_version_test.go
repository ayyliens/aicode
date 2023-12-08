package u_test

import (
	"_/go/u"
	"github.com/mitranim/gg"
	"testing"

	"github.com/mitranim/gg/gtest"
)

func Test_FileIndex_encoding(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Str(u.Version{0}, `0`)
	gtest.Str(u.Version{0, 1}, `0.1`)
	gtest.Str(u.Version{1}, `1`)
	gtest.Str(u.Version{1, 0}, `1.0`)
	gtest.Str(u.Version{2}, `2`)
	gtest.Str(u.Version{2, 1, 1}, `2.1.1`)
	gtest.Str(u.Version{9}, `9`)
	gtest.Str(u.Version{10}, `10`)
	gtest.Str(u.Version{19}, `19`)
	gtest.Str(u.Version{20}, `20`)
	gtest.Str(u.Version{99}, `99`)
	gtest.Str(u.Version{100}, `100`)
	gtest.Str(u.Version{999}, `999`)
	gtest.Str(u.Version{1000}, `1000`)
	gtest.Str(u.Version{9999}, `9999`)
}

func Test_FileIndex_decoding(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to decode "blah" into type uint16: strconv.ParseUint: parsing "blah": invalid syntax`,
		func() { gg.ParseTo[u.Version](`blah`) },
	)

	test := func(src string, exp u.Version) {
		gtest.True(gg.ParseTo[u.Version](src).Equal(exp))
	}

	test(`0.0.0.0`, u.Version{0, 0, 0, 0})
	test(`0.0.0.0`, u.Version{0})
	test(`0.0.0.1`, u.Version{0, 0, 0, 1})
	test(`1`, u.Version{1})
	test(`1`, u.Version{1, 0, 0})
	test(`2.0.0.0`, u.Version{2})
	test(`10`, u.Version{10, 0})
	test(`1.10.9.3`, u.Version{1, 10, 9, 3})
}
