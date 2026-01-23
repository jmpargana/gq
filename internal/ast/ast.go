package ast

import (
	u "github.com/jmpargana/gq/internal/utils"
)

func Transform(o any, n u.Node) any {
	prev := o

	// prevent doubled indexing when creating dict
	if n.Value.Kind == u.ASSIGN {
		return prev
	}

	for _, c := range n.Children {
		prev = Transform(prev, c)
	}

	switch n.Value.Kind {
	case u.IDX:
		prev = index(prev, n.Value)
	case u.INDEXSTART:
		prev = []any{prev}
	case u.DICTSTART:
		prev = dict(prev, n)
	}

	return prev
}

func dict(o any, n u.Node) any {
	m := map[string]any{}
	for _, c := range n.Children {
		m[c.Value.Ident] = Transform(o, c.Children[0])
	}
	return m
}

// FIXME: add error handling
func index(o any, c u.Cmd) any {
	prev := o
	for _, f := range c.Fields {
		if f.Kind == u.IDX {
			l := prev.([]any)
			prev = l[f.Idx]
		}
		if f.Kind == u.FIELD {
			m := prev.(map[string]any)
			prev = m[f.Name]
		}
	}
	return prev
}
