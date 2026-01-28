package ast

import (
	u "github.com/jmpargana/gq/internal/utils"
)

func TransformStream(s u.Stream, n u.Node) u.Stream {
	prev := s

	if n.Value.Kind == u.PIPE {

		// TODO: test chained pipes vs multiple children
		left := TransformStream(prev, n.Children[0])
		rightNode := n.Children[1]

		next := u.NewStream()

		for _, it := range left.O {
			out := TransformStream(u.NewSingleStream(it), rightNode)
			next.O = append(next.O, out.O...)
		}

		return next
	}

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
		prev = u.NewSingleStream(arr)
	case u.DICTSTART:
		prev = dictStream(prev, n)
	}

	return prev
}

func indexStream(s u.Stream, c u.Cmd) u.Stream {
	nextS := u.NewStream()
	for _, o := range s.O {
		prevs := []any{o}

		for _, f := range c.Fields {
			var newPrevs []any

			for _, prev := range prevs {
				switch f.Kind {
				case u.IDX:
					l := prev.([]any)
					newPrevs = append(newPrevs, l[f.Idx])
				case u.ROOT:
					newPrevs = append(newPrevs, prev)
				case u.FIELD:
					m := prev.(map[string]any)
					newPrevs = append(newPrevs, m[f.Name])
				case u.ARRAY:
					switch prev := prev.(type) {
					case []any:
						l := prev
						newPrevs = append(newPrevs, l...)
					case map[string]any:
						for _, v := range prev {
							newPrevs = append(newPrevs, v)
						}
					}
				}
			}

			prevs = newPrevs
		}

		nextS.O = append(nextS.O, prevs...)
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
