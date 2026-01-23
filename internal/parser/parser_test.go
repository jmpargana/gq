package parser

import (
	"testing"

	u "github.com/jmpargana/gq/internal/utils"
)

func TestBuildTree(t *testing.T) {
	testCases := []struct {
		desc string
		cmds []u.Cmd
		pgr  u.Node
	}{
		{
			desc: "single index",
			cmds: []u.Cmd{
				{Kind: u.IDX},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.IDX}},
		},
		{
			desc: "single depth",
			cmds: []u.Cmd{
				{Kind: u.INDEXSTART},
				{Kind: u.IDX},
				{Kind: u.INDEXEND},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX}}}},
		},
		{
			desc: "multiple depth",
			cmds: []u.Cmd{
				{Kind: u.INDEXSTART},
				{Kind: u.INDEXSTART},
				{Kind: u.IDX},
				{Kind: u.INDEXEND},
				{Kind: u.INDEXEND},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{
				{Value: u.Cmd{Kind: u.IDX}}}}}},
		},
		{
			desc: "u.PIPE",
			cmds: []u.Cmd{
				{Kind: u.IDX},
				{Kind: u.PIPE},
				{Kind: u.IDX},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.PIPE}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX}}, {Value: u.Cmd{Kind: u.IDX}}}},
		},
		{
			desc: "u.PIPE + depth",
			cmds: []u.Cmd{
				{Kind: u.INDEXSTART},
				{Kind: u.IDX},
				{Kind: u.INDEXEND},
				{Kind: u.PIPE},
				{Kind: u.IDX},
			},
			pgr: u.Node{Value: u.Cmd{Kind: u.PIPE}, Children: []u.Node{{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX}}}}, {Value: u.Cmd{Kind: u.IDX}}}},
		},
		{
			desc: "generates proper hierarchy",
			cmds: []u.Cmd{
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
			cmds: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN},
				{Kind: u.IDX},
				{Kind: u.DICTEND},
				{Kind: u.PIPE},
				{Kind: u.IDX},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.DICTSTART},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.ASSIGN}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX}}}},
						},
					},
					{Value: u.Cmd{Kind: u.IDX}},
				}},
		},
		{
			desc: "dict, multiple Values",
			cmds: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN},
				{Kind: u.IDX},
				{Kind: u.COMMA},
				{Kind: u.ASSIGN},
				{Kind: u.IDX},
				{Kind: u.COMMA},
				{Kind: u.ASSIGN},
				{Kind: u.IDX},
				{Kind: u.DICTEND},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.DICTSTART},
				Children: []u.Node{
					{Value: u.Cmd{Kind: u.ASSIGN}, Children: []u.Node{
						{Value: u.Cmd{Kind: u.IDX}},
					}},
					{Value: u.Cmd{Kind: u.ASSIGN}, Children: []u.Node{
						{Value: u.Cmd{Kind: u.IDX}},
					}},
					{Value: u.Cmd{Kind: u.ASSIGN}, Children: []u.Node{
						{Value: u.Cmd{Kind: u.IDX}},
					}},
				}},
		},
		{
			desc: "dict, array, u.PIPE",
			cmds: []u.Cmd{
				{Kind: u.DICTSTART},
				{Kind: u.ASSIGN},
				{Kind: u.IDX},
				{Kind: u.COMMA},
				{Kind: u.ASSIGN},
				{Kind: u.INDEXSTART},
				{Kind: u.IDX},
				{Kind: u.INDEXEND},
				{Kind: u.COMMA},
				{Kind: u.ASSIGN},
				{Kind: u.INDEXSTART},
				{Kind: u.IDX},
				{Kind: u.PIPE},
				{Kind: u.IDX},
				{Kind: u.INDEXEND},
				{Kind: u.DICTEND},
				{Kind: u.PIPE},
				{Kind: u.IDX},
			},
			pgr: u.Node{
				Value: u.Cmd{Kind: u.PIPE},
				Children: []u.Node{
					{
						Value: u.Cmd{Kind: u.DICTSTART},
						Children: []u.Node{
							{Value: u.Cmd{Kind: u.ASSIGN}, Children: []u.Node{
								{Value: u.Cmd{Kind: u.IDX}},
							}},
							{Value: u.Cmd{Kind: u.ASSIGN}, Children: []u.Node{
								{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{{Value: u.Cmd{Kind: u.IDX}}}},
							}},
							{Value: u.Cmd{Kind: u.ASSIGN}, Children: []u.Node{
								{
									Value: u.Cmd{Kind: u.INDEXSTART},
									Children: []u.Node{
										{Value: u.Cmd{Kind: u.PIPE}, Children: []u.Node{
											{Value: u.Cmd{Kind: u.IDX}},
											{Value: u.Cmd{Kind: u.IDX}},
										}},
									},
								},
							}},
						}},
					{Value: u.Cmd{Kind: u.IDX}},
				},
			},
		},
		// TODO: multiple chained u.PIPEs
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := NewParser(tC.cmds)
			got := p.ParseExpr()
			if !got.IsEqual(tC.pgr) {
				t.Fatalf("\nexpected: %v\ngot: %v\n", tC.pgr, got)
			}
		})
	}
}
