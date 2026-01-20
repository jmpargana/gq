package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestGQParser(t *testing.T) {
	testCases := []struct {
		desc string
		str  string
		pgr  []cmd
	}{
		{
			desc: "root",
			str:  ".",
			pgr: []cmd{
				{kind: idx, fields: []idxField{{kind: field, name: ""}}},
			},
		},
		{
			desc: "fieldA",
			str:  ".fieldA",
			pgr: []cmd{
				{kind: idx, fields: []idxField{{kind: field, name: "fieldA"}}},
			},
		},
		{
			desc: "fieldA under quotes",
			str:  ".\"fieldA\"",
			pgr: []cmd{
				{kind: idx, fields: []idxField{{kind: field, name: "fieldA"}}},
			},
		},
		{
			desc: "fieldA inside square brackets",
			str:  ".[fieldA]",
			pgr: []cmd{
				{kind: idx, fields: []idxField{{kind: field, name: "fieldA"}}},
			},
		},
		{
			desc: "index 0",
			str:  ".[0]",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: idx, idx: 0},
				}},
			},
		},
		{
			desc: "index 12",
			str:  ".[12]",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: idx, idx: 12},
				}},
			},
		},
		{
			desc: "fieldA fieldB nested",
			str:  ".fieldA.fieldB",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: field, name: "fieldA"},
					{kind: field, name: "fieldB"},
				}},
			},
		},
		{
			desc: "fieldA 0 nested",
			str:  ".fieldA.0",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: field, name: "fieldA"},
					{kind: idx, idx: 0},
				}},
			},
		},
		{
			desc: "fieldA 0 nested 2",
			str:  ".fieldA[0]",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: field, name: "fieldA"},
					{kind: idx, idx: 0},
				}},
			},
		},
		{
			desc: "fieldA 1 2 nested 3",
			str:  ".fieldA[1][2]",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: field, name: "fieldA"},
					{kind: idx, idx: 1},
					{kind: idx, idx: 2},
				}},
			},
		},
		{
			desc: "complex",
			str:  ".\"a\"[1][\"b\"][1]",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: field, name: "a"},
					{kind: idx, idx: 1},
					{kind: field, name: "b"},
					{kind: idx, idx: 1},
				}},
			},
		},
		{
			desc: "double",
			str:  ".[5][0]",
			pgr: []cmd{
				{kind: idx, fields: []idxField{
					{kind: idx, idx: 5},
					{kind: idx, idx: 0},
				}},
			},
		},
		{
			desc: "pipe",
			str:  ". | .",
			pgr: []cmd{
				{kind: idx, fields: []idxField{{kind: field, name: ""}}},
				{kind: pipe},
				{kind: idx, fields: []idxField{{kind: field, name: ""}}},
			},
		},
		{
			desc: "pipe with indexing",
			str:  ".a[5] | .[1].b",
			pgr: []cmd{
				{kind: idx, fields: []idxField{{kind: field, name: "a"}, {kind: idx, idx: 5}}},
				{kind: pipe},
				{kind: idx, fields: []idxField{{kind: idx, idx: 1}, {kind: field, name: "b"}}},
			},
		},
		{
			desc: "multiple pipes",
			str:  `."a"[5] | .5 | .b`,
			pgr: []cmd{
				{kind: idx, fields: []idxField{{kind: field, name: "a"}, {kind: idx, idx: 5}}},
				{kind: pipe},
				{kind: idx, fields: []idxField{{kind: idx, idx: 5}}},
				{kind: pipe},
				{kind: idx, fields: []idxField{{kind: field, name: "b"}}},
			},
		},
		{
			desc: "wrap in array",
			str:  `[."a"]`,
			pgr: []cmd{
				{kind: mapIndexStart},
				{kind: idx, fields: []idxField{{kind: field, name: "a"}}},
				{kind: mapIndexEnd},
			},
		},
		// {
		// 	desc: "multiple wrap with pipe",
		// 	str:  `[[."a" | .b] | .[0]]`,
		// 	pgr: []cmd{
		// 		{kind: mapIndexStart},
		// 		{kind: mapIndexStart},
		// 		{kind: idx, fields: []idxField{{kind: field, name: "a"}}},
		// 		{kind: pipe},
		// 		{kind: idx, fields: []idxField{{kind: field, name: "b"}}},
		// 		{kind: mapIndexEnd},
		// 		{kind: pipe},
		// 		{kind: idx, fields: []idxField{{kind: idx, idx: 0}}},
		// 		{kind: mapIndexEnd},
		// 	},
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := parseGQ(tC.str)

			if len(got) != len(tC.pgr) {
				t.Fatalf("didn't get right length")
			}

			for idx, it := range got {
				if !it.isEqual(tC.pgr[idx]) {
					t.Fatalf("not equal\nexpected: %v\ngot: %v", tC.pgr[idx], it)
				}
			}
		})
	}
}

func TestTranform(t *testing.T) {
	testCases := []struct {
		desc string
		a    string
		b    any
		pgr  []cmd
	}{
		{
			desc: "list index",
			a:    "[5, 10, 15, 20]",
			b:    int64(15),
			pgr:  []cmd{{kind: idx, fields: []idxField{{kind: idx, idx: 2}}}},
		},
		{
			desc: "list index",
			a:    "[5, 10, [21, 22], 20]",
			b:    int64(21),
			pgr:  []cmd{{kind: idx, fields: []idxField{{kind: idx, idx: 2}, {kind: idx, idx: 0}}}},
		},
		{
			desc: "map index",
			a:    "{\"a\": \"A\"}",
			b:    "A",
			pgr:  []cmd{{kind: idx, fields: []idxField{{kind: field, name: "a"}}}},
		},
		{
			desc: "nested list and map",
			a:    "{\"a\": [1, {\"b\": [2, 3]}]}",
			b:    int64(3),
			pgr:  []cmd{{kind: idx, fields: []idxField{{kind: field, name: "a"}, {kind: idx, idx: 1}, {kind: field, name: "b"}, {kind: idx, idx: 1}}}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a := parseObject(bufio.NewReader(strings.NewReader(tC.a)))
			got := transform(a, tC.pgr)
			if tC.b != got {
				t.Fatalf("not equal:\ngot: %v\nwanted: %v", got, tC.b)
			}
		})
	}
}

func (n node) isEqual(o node) bool {
	if !n.value.isEqual(o.value) || len(n.children) != len(o.children) {
		return false
	}
	for i := range n.children {
		if !n.children[i].isEqual(o.children[i]) {
			return false
		}
	}
	return true
}

func TestBuildTree(t *testing.T) {
	testCases := []struct {
		desc string
		cmds []cmd
		pgr  node
	}{
		{
			desc: "single index",
			cmds: []cmd{
				{kind: idx},
			},
			pgr: node{value: cmd{kind: idx}},
		},
		{
			desc: "single depth",
			cmds: []cmd{
				{kind: mapIndexStart},
				{kind: idx},
				{kind: mapIndexEnd},
			},
			pgr: node{value: cmd{kind: mapIndexStart}, children: []node{{value: cmd{kind: idx}}}},
		},
		{
			desc: "multiple depth",
			cmds: []cmd{
				{kind: mapIndexStart},
				{kind: mapIndexStart},
				{kind: idx},
				{kind: mapIndexEnd},
				{kind: mapIndexEnd},
			},
			pgr: node{value: cmd{kind: mapIndexStart}, children: []node{{value: cmd{kind: mapIndexStart}, children: []node{
				{value: cmd{kind: idx}}}}}},
		},
		{
			desc: "pipe",
			cmds: []cmd{
				{kind: idx},
				{kind: pipe},
				{kind: idx},
			},
			pgr: node{value: cmd{kind: pipe}, children: []node{{value: cmd{kind: idx}}, {value: cmd{kind: idx}}}},
		},
		{
			desc: "pipe + depth",
			cmds: []cmd{
				{kind: mapIndexStart},
				{kind: idx},
				{kind: mapIndexEnd},
				{kind: pipe},
				{kind: idx},
			},
			pgr: node{value: cmd{kind: pipe}, children: []node{{value: cmd{kind: mapIndexStart}, children: []node{{value: cmd{kind: idx}}}}, {value: cmd{kind: idx}}}},
		},
		// {
		// 	desc: "generates proper hierarchy",
		// 	cmds: []cmd{
		// 		{kind: mapIndexStart},
		// 		{kind: mapIndexStart},
		// 		{kind: idx, fields: []idxField{{kind: field, name: "a"}}},
		// 		{kind: pipe},
		// 		{kind: idx, fields: []idxField{{kind: field, name: "b"}}},
		// 		{kind: mapIndexEnd},
		// 		{kind: pipe},
		// 		{kind: idx, fields: []idxField{{kind: idx, idx: 0}}},
		// 		{kind: mapIndexEnd},
		// 	},
		// 	pgr: node{
		// 		value: cmd{kind: mapIndexStart},
		// 		children: []node{
		// 			{value: cmd{kind: pipe},
		// 				children: []node{
		// 					{value: cmd{kind: mapIndexStart},
		// 						children: []node{
		// 							{value: cmd{kind: pipe},
		// 								children: []node{
		// 									{value: cmd{kind: idx, fields: []idxField{{kind: field, name: "a"}}}},
		// 									{
		// 										value: cmd{kind: idx, fields: []idxField{{kind: field, name: "b"}}},
		// 									},
		// 								},
		// 							},
		// 							{
		// 								value: cmd{kind: idx, fields: []idxField{{kind: idx, idx: 0}}},
		// 							},
		// 						}}}},
		// 		},
		// 	},
		// }
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := buildTree(tC.cmds)
			if !got.isEqual(tC.pgr) {
				t.Fatalf("\nexpected: %v\ngot: %v\n", tC.pgr, got)
			}
		})
	}
}
