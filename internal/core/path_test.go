package core_test

import (
	"reflect"
	"testing"

	"github.com/ZentriaMC/apply-config/internal/core"
)

func TestPathParsing(t *testing.T) {
	type testCase struct {
		path     string
		expected []core.PathElement
	}

	testCases := []testCase{
		{`one.two.three`, core.PathElements(`one`, `two`, `three`)},
		{`one."two".three`, core.PathElements(`one`, `two`, `three`)},
		{`one."two"."three"`, core.PathElements(`one`, `two`, `three`)},
		{`one.\"two\".three`, core.PathElements(`one`, `"two"`, `three`)},
		{`one."two\"".three`, core.PathElements(`one`, `two"`, `three`)},
		{`one."two"`, core.PathElements(`one`, `two`)},
		{`one[48].two.\[three`, core.PathElements(`one`, 48, `two`, `[three`)},
	}

	for _, testCase := range testCases {
		path := testCase.path
		expected := testCase.expected

		res, err := core.ProcessPath(path)
		if err != nil {
			t.Fatal(err)
			return
		}

		t.Logf("paths: '%s' -> %+v", path, res)
		if !reflect.DeepEqual(expected, res) {
			t.Errorf("mismatch, %+v != %+v", expected, res)
		}
	}
}
