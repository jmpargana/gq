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
			got := ParseObject(bufio.NewReader(strings.NewReader(tC.s)))
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
			got := ParseObject(bufio.NewReader(strings.NewReader(tC.s)))
			if !reflect.DeepEqual(got, tC.arr) {
				t.Fatalf("failed comparison\ngot: %v\nexpected: %v\n", got, tC.arr)
			}
		})
	}
}
