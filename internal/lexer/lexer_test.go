package lexer

import (
	"testing"

	u "github.com/jmpargana/gq/internal/utils"
)

func TestGQParser(t *testing.T) {
	testCases := []struct {
		desc string
		str  string
		pgr  []u.Cmd
	}{
		{
			desc: "root",
			str:  ".",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: ""}}},
			},
		},
		{
			desc: "FieldA",
			str:  ".FieldA",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "FieldA"}}},
			},
		},
		{
			desc: "FieldA under quotes",
			str:  ".\"FieldA\"",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "FieldA"}}},
			},
		},
		{
			desc: "FieldA inside square brackets",
			str:  ".[FieldA]",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "FieldA"}}},
			},
		},
		{
			desc: "index 0",
			str:  ".[0]",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.IDX, Idx: 0},
				}},
			},
		},
		{
			desc: "index 12",
			str:  ".[12]",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.IDX, Idx: 12},
				}},
			},
		},
		{
			desc: "FieldA fieldB nested",
			str:  ".FieldA.FieldB",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.FIELD, Name: "FieldA"},
					{Kind: u.FIELD, Name: "FieldB"},
				}},
			},
		},
		{
			desc: "FieldA 0 nested",
			str:  ".FieldA.0",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.FIELD, Name: "FieldA"},
					{Kind: u.IDX, Idx: 0},
				}},
			},
		},
		{
			desc: "FieldA 0 nested 2",
			str:  ".FieldA[0]",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.FIELD, Name: "FieldA"},
					{Kind: u.IDX, Idx: 0},
				}},
			},
		},
		{
			desc: "FieldA 1 2 nested 3",
			str:  ".FieldA[1][2]",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.FIELD, Name: "FieldA"},
					{Kind: u.IDX, Idx: 1},
					{Kind: u.IDX, Idx: 2},
				}},
			},
		},
		{
			desc: "complex",
			str:  ".\"a\"[1][\"b\"][1]",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.FIELD, Name: "a"},
					{Kind: u.IDX, Idx: 1},
					{Kind: u.FIELD, Name: "b"},
					{Kind: u.IDX, Idx: 1},
				}},
			},
		},
		{
			desc: "double",
			str:  ".[5][0]",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{
					{Kind: u.IDX, Idx: 5},
					{Kind: u.IDX, Idx: 0},
				}},
			},
		},
		{
			desc: "u.PIPE",
			str:  ". | .",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: ""}}},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: ""}}},
			},
		},
		{
			desc: "u.PIPE with indexing",
			str:  ".a[5] | .[1].b",
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.IDX, Idx: 5}}},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 1}, {Kind: u.FIELD, Name: "b"}}},
			},
		},
		{
			desc: "multiple u.PIPEs",
			str:  `."a"[5] | .5 | .b`,
			pgr: []u.Cmd{
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.IDX, Idx: 5}}},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 5}}},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
			},
		},
		{
			desc: "wrap in array",
			str:  `[."a"]`,
			pgr: []u.Cmd{
				{Kind: u.INDEXSTART},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}},
				{Kind: u.INDEXEND},
			},
		},
		{
			desc: "multiple wrap with u.PIPE",
			str:  `[[."a" | .b] | .[0]]`,
			pgr: []u.Cmd{
				{Kind: u.INDEXSTART},
				{Kind: u.INDEXSTART},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.INDEXEND},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 0}}},
				{Kind: u.INDEXEND},
			},
		},
		{
			desc: "dict single value",
			str:  `{a: .b}`,
			pgr: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.DICTEND},
			},
		},
		{
			desc: "multiple dicts",
			str:  `{a: .b, b: .c}`,
			pgr: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.COMMA},
				{Kind: u.ASSIGN, Ident: "b"},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "c"}}},
				{Kind: u.DICTEND},
			},
		},
		{
			desc: "multiple dicts",
			str:  `{a: [.b], b: .c}`,
			pgr: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.INDEXSTART},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.INDEXEND},
				{Kind: u.COMMA},
				{Kind: u.ASSIGN, Ident: "b"},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "c"}}},
				{Kind: u.DICTEND},
			},
		},
		{
			desc: "dict array",
			str:  `{a: [.b]}`,
			pgr: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.INDEXSTART},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.INDEXEND},
				{Kind: u.DICTEND},
			},
		},
		{
			desc: "dict array",
			str:  `{a: [.b | .c]}`,
			pgr: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.INDEXSTART},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "c"}}},
				{Kind: u.INDEXEND},
				{Kind: u.DICTEND},
			},
		},
		{
			desc: "dict array",
			str:  `[{a: .b}]`,
			pgr: []u.Cmd{
				{Kind: u.INDEXSTART},
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.DICTEND},
				{Kind: u.INDEXEND},
			},
		},
		{
			desc: "dict array filter",
			str:  `{a: .b} | .c`,
			pgr: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.DICTEND},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "c"}}},
			},
		},
		{
			desc: "dict array filter array",
			str:  `[{a: .b} | .c]`,
			pgr: []u.Cmd{
				{Kind: u.INDEXSTART},
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN, Ident: "a"},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
				{Kind: u.DICTEND},
				{Kind: u.PIPE},
				{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "c"}}},
				{Kind: u.INDEXEND},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := Lex(tC.str)

			if len(got) != len(tC.pgr) {
				t.Fatalf("didn't get right length")
			}

			for idx, it := range got {
				if !it.IsEqual(tC.pgr[idx]) {
					t.Fatalf("not equal\nexpected: %v\ngot: %v", tC.pgr[idx], it)
				}
			}
		})
	}
}
