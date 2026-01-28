package parser

import (
	"reflect"
	"testing"

	l "github.com/jmpargana/gq/internal/lexer"
	u "github.com/jmpargana/gq/internal/utils"
)

func TestBuildTree(t *testing.T) {
	testCases := []struct {
		desc string
		cmds []l.Token
		pgr  u.Node
	}{
		{
			desc: "single root index",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
		},
		{
			desc: "single number",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.NUMBER, Value: "5"},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 5}}}},
		},
		{
			desc: "single ident",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.IDENT, Value: "a"},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}}},
		},
		{
			desc: "single string",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.STRING, Value: "a"},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}}},
		},
		{
			desc: "single string inside braces",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.STRING, Value: "a"},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}}},
		},
		{
			desc: "single ident inside braces",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.IDENT, Value: "a"},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}}},
		},
		{
			desc: "single iterator",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
		},
		{
			desc: "double number",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.NUMBER, Value: "5"},
				{Kind: l.RBRACE},
				{Kind: l.LBRACE},
				{Kind: l.NUMBER, Value: "4"},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 5}, {Kind: u.IDX, Idx: 4}}}},
		},
		{
			desc: "double number with dot",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.NUMBER, Value: "5"},
				{Kind: l.RBRACE},
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.NUMBER, Value: "4"},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 5}, {Kind: u.IDX, Idx: 4}}}},
		},
		{
			desc: "ident + ident",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.IDENT, Value: "a"},
				{Kind: l.DOT},
				{Kind: l.IDENT, Value: "b"},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.FIELD, Name: "b"}}}},
		},
		{
			desc: "string + string in braces",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.STRING, Value: "a"},
				{Kind: l.LBRACE},
				{Kind: l.STRING, Value: "b"},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.FIELD, Name: "b"}}}},
		},
		{
			desc: "string + iterator",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.STRING, Value: "a"},
				{Kind: l.LBRACE},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}, {Kind: u.ARRAY}}}},
		},
		{
			desc: "single depth",
			cmds: []l.Token{
				{Kind: l.LBRACE},
				{Kind: l.DOT},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}}}},
		},
		{
			desc: "multiple depth",
			cmds: []l.Token{
				{Kind: l.LBRACE},
				{Kind: l.LBRACE},
				{Kind: l.DOT},
				{Kind: l.RBRACE},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{
				{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}}}}}},
		},
		{
			desc: "u.PIPE",
			cmds: []l.Token{
				{Kind: l.DOT},
				{Kind: l.PIPE},
				{Kind: l.DOT},
				{Kind: l.EOF},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.PIPE}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}}, {Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}}}},
		},
		{
			desc: "u.PIPE + depth",
			cmds: []l.Token{
				{Kind: l.LBRACE},
				{Kind: l.DOT},
				{Kind: l.RBRACE},
				{Kind: l.PIPE},
				{Kind: l.DOT},
				{Kind: l.EOF},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.INDEXSTART},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
						},
					},
					{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
				},
			},
		},
		{
			desc: "generates proper hierarchy",
			cmds: []l.Token{
				{Kind: l.LBRACE},
				{Kind: l.LBRACE},
				{Kind: l.DOT},
				{Kind: l.STRING, Value: "a"},
				{Kind: l.PIPE},
				{Kind: l.DOT},
				{Kind: l.STRING, Value: "b"},
				{Kind: l.RBRACE},
				{Kind: l.PIPE},
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.NUMBER, Value: "0"},
				{Kind: l.RBRACE},
				{Kind: l.RBRACE},
				{Kind: l.EOF},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.INDEXSTART},
				Children: []u.Node{
					{Value: u.Cmd{Kind: u.PIPE},
						Children: []u.Node{
							{
								Value: u.Cmd{Kind: u.INDEXSTART},
								Children: []u.Node{
									{Value: u.Cmd{Kind: u.PIPE},
										Children: []u.Node{
											{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "a"}}}},
											{
												Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}},
											},
										},
									}},
							},
							{
								Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.IDX, Idx: 0}}},
							},
						}},
				},
			},
		},
		{
			desc: "u.PIPE + dict",
			cmds: []l.Token{
				{Kind: l.LBRACKET},
				{Kind: l.IDENT, Value: "a"},
				{Kind: l.COLON},
				{Kind: l.DOT},
				{Kind: l.IDENT, Value: "b"},
				{Kind: l.RBRACKET},
				{Kind: l.PIPE},
				{Kind: l.DOT},
				{Kind: l.EOF},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.DICTSTART},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.ASSIGN, Ident: "a"}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "b"}}}}}},
						},
					},
					{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
				}},
		},
		{
			desc: "dict, multiple Values",
			cmds: []l.Token{
				{Kind: l.LBRACKET},
				{Kind: l.IDENT, Value: "a"},
				{Kind: l.COLON},
				{Kind: l.DOT},
				{Kind: l.COMMA},
				{Kind: l.IDENT, Value: "b"},
				{Kind: l.COLON},
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.RBRACE},
				{Kind: l.COMMA},
				{Kind: l.IDENT, Value: "c"},
				{Kind: l.COLON},
				{Kind: l.DOT},
				{Kind: l.LBRACE},
				{Kind: l.IDENT, Value: "d"},
				{Kind: l.RBRACE},
				{Kind: l.RBRACKET},
				{Kind: l.EOF},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.DICTSTART},
				Children: []u.Node{
					{Value: u.Cmd{Kind: u.ASSIGN, Ident: "a"}, Children: []u.Node{
						{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
					}},
					{Value: u.Cmd{Kind: u.ASSIGN, Ident: "b"}, Children: []u.Node{
						{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ARRAY}}}},
					}},
					{Value: u.Cmd{Kind: u.ASSIGN, Ident: "c"}, Children: []u.Node{
						{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.FIELD, Name: "d"}}}},
					}},
				}},
		},
		{
			desc: "dict, array, u.PIPE",
			cmds: []l.Token{
				{Kind: l.LBRACKET},
				{Kind: l.IDENT, Value: "a"},
				{Kind: l.COLON},
				{Kind: l.DOT},
				{Kind: l.COMMA},
				{Kind: l.IDENT, Value: "b"},
				{Kind: l.COLON},
				{Kind: l.LBRACE},
				{Kind: l.DOT},
				{Kind: l.RBRACE},
				{Kind: l.COMMA},
				{Kind: l.IDENT, Value: "c"},
				{Kind: l.COLON},
				{Kind: l.LBRACE},
				{Kind: l.DOT},
				{Kind: l.PIPE},
				{Kind: l.DOT},
				{Kind: l.RBRACE},
				{Kind: l.RBRACKET},
				{Kind: l.PIPE},
				{Kind: l.DOT},
				{Kind: l.EOF},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.DICTSTART},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.ASSIGN, Ident: "a"}, Children: []u.Node{
								{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
							}},
							{Value: u.Cmd{Kind: u.ASSIGN, Ident: "b"}, Children: []u.Node{
								{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}}}},
							}},
							{Value: u.Cmd{Kind: u.ASSIGN, Ident: "c"}, Children: []u.Node{
								{
									Value: u.Cmd{Kind: u.INDEXSTART},
									Children: []u.Node{
										{Value: u.Cmd{Kind: u.PIPE}, Children: []u.Node{
											{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
											{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
										}},
									},
								},
							}},
						}},
					{Value: u.Cmd{Kind: u.IDX, Fields: []u.IdxField{{Kind: u.ROOT}}}},
				},
			},
		},
		// TODO: multiple chained u.PIPEs
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := NewParser(tC.cmds)
			got := p.ParseExpr()
			if !reflect.DeepEqual(tC.pgr, got) {
				t.Fatalf("\nexpected:\n\t%v\ngot:\n\t%v\n", tC.pgr, got)
			}
		})
	}
}
