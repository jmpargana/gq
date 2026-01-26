package ast

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	json "github.com/jmpargana/gq/internal/gqjson"
	u "github.com/jmpargana/gq/internal/utils"
)

func TestTranform(t *testing.T) {
	testCases := []struct {
		desc string
		a    string
		b    any
		pgr  u.Node
	}{
		{
			desc: "list index",
			a:    "[5, 10, 15, 20]",
			b:    int64(15),
			pgr: u.Node{
				Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 2}}},
			},
		},
		{
			desc: "root is same",
			a:    "[5, 10, 15, 20]",
			b:    []interface{}{int64(5), int64(10), int64(15), int64(20)},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}},
			},
		},
		{
			desc: "list index",
			a:    "[5, 10, [21, 22], 20]",
			b:    int64(21),
			pgr: u.Node{
				Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 2}, {Kind: u.IDX, Idx: 0}}},
			},
		},
		{
			desc: "map index",
			a:    "{\"a\": \"A\"}",
			b:    "A",
			pgr: u.Node{
				Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}},
			},
		},
		{
			desc: "nested list and map",
			a:    "{\"a\": [1, {\"b\": [2, 3]}]}",
			b:    int64(3),
			pgr: u.Node{
				Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.IDX, Idx: 1}, {Kind: u.FIELD, Name: "b"}, {Kind: u.IDX, Idx: 1}}},
			},
		},
		{
			desc: "nested list and map",
			a:    "{\"a\": [1, {\"b\": [2, 3]}]}",
			b:    []interface{}{int64(3)},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.INDEXSTART},
				Children: []u.Node{{
					Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.IDX, Idx: 1}, {Kind: u.FIELD, Name: "b"}, {Kind: u.IDX, Idx: 1}}},
				}},
			},
		},
		{
			desc: "nested list and map",
			a:    "{\"a\": [1, {\"b\": [2, 3]}]}",
			b:    int64(3),
			pgr: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.INDEXSTART},
						Children: []u.Node{{
							Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.IDX, Idx: 1}, {Kind: u.FIELD, Name: "b"}, {Kind: u.IDX, Idx: 1}}},
						}},
					},
					{
						Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 0}}},
					},
				},
			},
		},
		{
			desc: "map with new key",
			a:    "{\"a\": [1, {\"b\": [2, 3]}]}",
			b:    map[string]interface{}{"b": int64(3)},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.DICTSTART},
				Children: []u.Node{{
					Value: u.Cmd{Kind: u.ASSIGN, Ident: "b"},
					Children: []u.Node{
						{
							Value: u.Cmd{Kind: u.PIPE},
							Children: []u.Node{
								{
									Value: u.Cmd{Kind: u.INDEXSTART},
									Children: []u.Node{{
										Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{
											{Kind: u.FIELD, Name: "a"},
											{Kind: u.IDX, Idx: 1},
											{Kind: u.FIELD, Name: "b"},
											{Kind: u.IDX, Idx: 1},
										},
										},
									}},
								},
								{
									Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 0}}},
								},
							},
						},
					},
				}},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			m := map[string]any{}

			m["a"] = "b"

			a := json.ParseObject(bufio.NewReader(strings.NewReader(tC.a)))
			got := Transform(a, tC.pgr)
			if !reflect.DeepEqual(tC.b, got) {
				t.Fatalf("not equal:\ngot: %v\nwanted: %v", got, tC.b)
			}
		})
	}
}
