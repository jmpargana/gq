package gqjson

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestArray(t *testing.T) {
	testCases := []struct {
		desc string
		s    string
		arr  []interface{}
	}{
		{
			desc: "multiple entities",
			s:    `["a", 3, 4.2, true, [1, 2], {"a": "b"}]`,
			arr:  []interface{}{"a", int64(3), float64(4.2), true, []interface{}{int64(1), int64(2)}, map[string]interface{}{"a": "b"}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := ParseObject(bufio.NewReader(strings.NewReader(tC.s)))
			if err != nil {
				t.Fatalf("expected no error, instead got: %v", err)
			}
			if !reflect.DeepEqual(got, tC.arr) {
				t.Fatalf("failed comparison\ngot: %v\nexpected: %v\n", got, tC.arr)
			}
		})
	}
}

func TestObject(t *testing.T) {
	testCases := []struct {
		desc string
		s    string
		arr  map[string]interface{}
	}{
		{
			desc: "multiple entities",
			s:    `{"a": "b", "b": 2, "c": true, "d": [1, 2], "e": {"a": 1}}`,
			arr: map[string]interface{}{
				"a": "b",
				"b": int64(2),
				"c": true,
				"d": []interface{}{int64(1), int64(2)},
				"e": map[string]interface{}{
					"a": int64(1),
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := ParseObject(bufio.NewReader(strings.NewReader(tC.s)))
			if err != nil {
				t.Fatalf("expected no error, instead got: %v", err)
			}
			if !reflect.DeepEqual(got, tC.arr) {
				t.Fatalf("failed comparison\ngot: %v\nexpected: %v\n", got, tC.arr)
			}
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	testCases := []struct {
		desc, s, err string
	}{
		{
			desc: "broken bool",
			s:    `{"a": ta}`,
			err:  "failed closing bool",
		},
		{
			desc: "broken int",
			s:    `{"a": 8a}`,
			err:  "failed parsing int",
		},
		{
			desc: "broken int",
			s:    `{"a": 8.a}`,
			err:  "failed parsing float",
		},
		{
			desc: "broken float",
			s:    `{"a": 8.a}`,
			err:  "failed parsing float",
		},
		{
			desc: "unclosed array",
			s:    `[1, 2`,
			err:  "unclosed array",
		},
		{
			desc: "unclosed dict",
			s:    `{`,
			err:  "invalid object",
		},
		{
			desc: "obj",
			s:    ``,
			err:  "invalid object",
		},
		{
			desc: "obj",
			s:    `{"a": [}`,
			err:  "unclosed array",
		},
		{
			desc: "string",
			s:    `{"a": "\\`,
			err:  "failed parsing string",
		},
		{
			desc: "string",
			s:    `{"a": "asldkj`,
			err:  "failed parsing string",
		},
		{
			desc: "unclosed nested",
			s:    `{"a": {"b": }`,
			err:  "invalid object",
		},
		{
			desc: "unclosed nested",
			s:    `[1, 2, {"a"]`,
			err:  "invalid object",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := ParseObject(bufio.NewReader(strings.NewReader(tC.s)))
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), tC.err) {
				t.Fatalf("expected %s\ngot: %s\n", tC.err, err)
			}
		})
	}
}
