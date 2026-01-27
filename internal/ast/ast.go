package ast

import (
	u "github.com/jmpargana/gq/internal/utils"
)

func TransformStream(s u.Stream, n u.Node) u.Stream {
	prev := s

	// prevent doubled indexing when creating dict
	if n.Value.Kind == u.ASSIGN {
		return prev
	}

	for _, c := range n.Children {
		prev = TransformStream(prev, c)
	}

	switch n.Value.Kind {
	case u.IDX:
		prev = indexStream(prev, n.Value)
	case u.INDEXSTART:
		arr := []any{}
		arr = append(arr, prev.O...)
		nextS := u.NewStream()
		nextS.O = append(nextS.O, arr)
		prev = nextS
	case u.DICTSTART:
		prev = dictStream(prev, n)
	}

	return prev
}

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

func indexStream(s u.Stream, c u.Cmd) u.Stream {
	nextS := u.NewStream()
	for _, o := range s.O {
		prev := o

		for _, f := range c.Fields {
			switch f.Kind {
			case u.IDX:
				l := prev.([]any)
				prev = l[f.Idx]
			case u.FIELD:
				m := prev.(map[string]any)
				prev = m[f.Name]
			case u.ARRAY:
				l := prev.([]any)
				nextS.O = append(nextS.O, l...)
				return nextS
			}
		}

		nextS.O = append(nextS.O, prev)
	}
	return nextS
}

func cloneMap(m map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func dictStream(o u.Stream, n u.Node) u.Stream {
	nextS := u.NewStream()
	for _, s := range o.O {
		// cartesian product
		partials := []map[string]any{
			{},
		}

		for _, c := range n.Children {

			innerS := u.NewStream()
			innerS.O = append(innerS.O, s)

			innerS = TransformStream(innerS, c.Children[0])

			var nextPartials []map[string]any

			for _, p := range partials {
				for _, in := range innerS.O {
					np := cloneMap(p)
					np[c.Value.Ident] = in
					nextPartials = append(nextPartials, np)
				}
			}

			partials = nextPartials

		}

		for _, p := range partials {
			nextS.O = append(nextS.O, p)
		}
	}
	return nextS
}
