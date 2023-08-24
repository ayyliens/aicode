package u_test

import (
	"_/go/u"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func Test_FileIndex_encoding(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Str(u.FileIndex(0), `0000`)
	gtest.Str(u.FileIndex(1), `0001`)
	gtest.Str(u.FileIndex(2), `0002`)
	gtest.Str(u.FileIndex(9), `0009`)
	gtest.Str(u.FileIndex(10), `0010`)
	gtest.Str(u.FileIndex(19), `0019`)
	gtest.Str(u.FileIndex(20), `0020`)
	gtest.Str(u.FileIndex(99), `0099`)
	gtest.Str(u.FileIndex(100), `0100`)
	gtest.Str(u.FileIndex(999), `0999`)
	gtest.Str(u.FileIndex(1000), `1000`)
	gtest.Str(u.FileIndex(9999), `9999`)

	gtest.PanicStr(
		`FileIndex 10000 overflows allowed digit count 4`,
		func() { u.FileIndex(10000).String() },
	)
}

func Test_FileIndex_decoding(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to decode "0" as *u.FileIndex: length mismatch: expected 4 digits, found 1 digits`,
		func() { gg.ParseTo[u.FileIndex](`0`) },
	)

	gtest.PanicStr(
		`unable to decode "00" as *u.FileIndex: length mismatch: expected 4 digits, found 2 digits`,
		func() { gg.ParseTo[u.FileIndex](`00`) },
	)

	gtest.PanicStr(
		`unable to decode "000" as *u.FileIndex: length mismatch: expected 4 digits, found 3 digits`,
		func() { gg.ParseTo[u.FileIndex](`000`) },
	)

	gtest.PanicStr(
		`unable to decode "00000" as *u.FileIndex: length mismatch: expected 4 digits, found 5 digits`,
		func() { gg.ParseTo[u.FileIndex](`00000`) },
	)

	gtest.PanicStr(
		`unable to decode "blah" into type uint16: strconv.ParseUint: parsing "blah": invalid syntax`,
		func() { gg.ParseTo[u.FileIndex](`blah`) },
	)

	test := func(src string, exp u.FileIndex) {
		gtest.Eq(
			gg.ParseTo[u.FileIndex](src),
			exp,
		)
	}

	test(`0000`, 0)
	test(`0001`, 1)
	test(`0002`, 2)
	test(`0010`, 10)
	test(`0012`, 12)
	test(`0120`, 120)
	test(`0123`, 123)
	test(`1230`, 1230)
	test(`1234`, 1234)
	test(`9999`, 9999)
}
