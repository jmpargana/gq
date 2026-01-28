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
			a := json.ParseObject(bufio.NewReader(strings.NewReader(tC.a)))
			s := u.NewSingleStream(a)
			got := TransformStream(s, tC.pgr)
			expected := u.NewSingleStream(tC.b)
			if !reflect.DeepEqual(expected, got) {
				t.Fatalf("not equal:\ngot: %v\nwanted: %v", got, expected)
			}
		})
	}
}

func TestTransformStream(t *testing.T) {
	testCases := []struct {
		desc    string
		start   string
		result  u.Stream
		program u.Node
	}{
		{
			desc:   "Array to stream",
			start:  `[1, 2, 3, 4]`,
			result: u.Stream{O: []any{int64(1), int64(2), int64(3), int64(4)}},
			program: u.Node{
				Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}},
			},
		},
		{
			desc:   "Array to stream to array",
			start:  `[1, 2, 3, 4]`,
			result: u.Stream{O: []any{[]interface{}{int64(1), int64(2), int64(3), int64(4)}}},
			program: u.Node{
				Value:    u.Cmd{Kind: u.INDEXSTART},
				Children: []u.Node{{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}}},
			},
		},
		{
			desc:   "piped index to array",
			start:  `[{"a": "b"}, {"a": "c"}]`,
			result: u.Stream{O: []any{[]interface{}{"b", "c"}}},
			program: u.Node{
				Value: u.Cmd{Kind: u.INDEXSTART},
				Children: []u.Node{{
					Value: u.Cmd{Kind: u.PIPE},
					Children: []u.Node{
						{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
						{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}}},
					},
				},
				},
			},
		},
		{
			// '.[] | {letter: .a}'
			desc:   "piped dict",
			start:  `[{"a": "b"}, {"a": "c"}]`,
			result: u.Stream{O: []any{map[string]any{"letter": "b"}, map[string]any{"letter": "c"}}},
			program: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}},
					},
					{
						Value: u.Cmd{Kind: u.DICTSTART},
						Children: []u.Node{
							{
								Value: u.Cmd{Kind: u.ASSIGN, Ident: "letter"},
								Children: []u.Node{
									{
										Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			// '{a: .[]}'
			desc:   "streamed dict",
			start:  `[[1], [2]]`,
			result: u.Stream{O: []any{map[string]any{"a": []interface{}{int64(1)}}, map[string]any{"a": []interface{}{int64(2)}}}},
			program: u.Node{
				Value: u.Cmd{Kind: u.DICTSTART},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.ASSIGN, Ident: "a"},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
						},
					},
				},
			},
		},
		// FIXME: not correct either
		{
			// '.[] | {a: .[]}'
			desc:   "streamed dict",
			start:  `[[1], [2]]`,
			result: u.Stream{O: []any{map[string]any{"a": int64(1)}, map[string]any{"a": int64(2)}}},
			program: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
					{
						Value: u.Cmd{Kind: u.DICTSTART},
						Children: []u.Node{
							{
								Value: u.Cmd{Kind: u.ASSIGN, Ident: "a"},
								Children: []u.Node{
									{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
								},
							},
						},
					},
				},
			},
		},
		{
			// '[.[]]'
			desc:   "array of iter",
			start:  `[1, 2]`,
			result: u.Stream{O: []any{[]any{int64(1), int64(2)}}},
			program: u.Node{
				Value: u.Cmd{Kind: u.INDEXSTART},
				Children: []u.Node{
					{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
				},
			},
		},
		{
			// '.[] | [.]'
			desc:   "piped array",
			start:  `[1, 2]`,
			result: u.Stream{O: []any{[]any{int64(1)}, []any{int64(2)}}},
			program: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
					{
						Value: u.Cmd{Kind: u.INDEXSTART},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
						},
					},
				},
			},
		},
		// FIXME: this stream should create cartesian product. It's not
		{
			// '.[] | {a: .[], b: .[]}'
			desc:  "streamed dict",
			start: `[[1], [2]]`,
			result: u.Stream{O: []any{
				map[string]any{"a": int64(1), "b": int64(1)},
				map[string]any{"a": int64(2), "b": int64(2)},
			}},
			program: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
					{
						Value: u.Cmd{Kind: u.DICTSTART},
						Children: []u.Node{
							{
								Value: u.Cmd{Kind: u.ASSIGN, Ident: "a"},
								Children: []u.Node{
									{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
								},
							},
							{
								Value: u.Cmd{Kind: u.ASSIGN, Ident: "b"},
								Children: []u.Node{
									{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
								},
							},
						},
					},
				},
			},
		},
		{
			// '{a: .[], b: .[]}'
			desc:  "streamed dict",
			start: `[[1], [2]]`,
			result: u.Stream{O: []any{
				map[string]any{"a": []interface{}{int64(1)}, "b": []interface{}{int64(1)}},
				map[string]any{"a": []interface{}{int64(1)}, "b": []interface{}{int64(2)}},
				map[string]any{"a": []interface{}{int64(2)}, "b": []interface{}{int64(1)}},
				map[string]any{"a": []interface{}{int64(2)}, "b": []interface{}{int64(2)}},
			}},
			program: u.Node{
				Value: u.Cmd{Kind: u.DICTSTART},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.ASSIGN, Ident: "a"},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
						},
					},
					{
						Value: u.Cmd{Kind: u.ASSIGN, Ident: "b"},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
						},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a := json.ParseObject(bufio.NewReader(strings.NewReader(tC.start)))
			s := u.NewStream()
			s.O = []any{a}
			got := TransformStream(s, tC.program)
			if !reflect.DeepEqual(tC.result, got) {
				t.Fatalf("not equal:\ngot: %v\nwanted: %v", got, tC.result)
			}
		})
	}
}
