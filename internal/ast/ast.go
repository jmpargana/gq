package ast

import (
	"github.com/jmpargana/gq/internal/stream"
	u "github.com/jmpargana/gq/internal/utils"
)

func TransformStream(s stream.Stream, n u.Node) stream.Stream {
	prev := s

	if n.Value.Kind == u.PIPE {

		// TODO: test chained pipes vs multiple children
		left := TransformStream(prev, n.Children[0])
		rightNode := n.Children[1]

		next := stream.New()

		for _, it := range left.O {
			out := TransformStream(stream.NewS(it), rightNode)
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
		prev = stream.NewS(arr)
	case u.DICTSTART:
		prev = dictStream(prev, n)
	}

	return prev
}

func indexStream(s stream.Stream, c u.Cmd) stream.Stream {
	nextS := stream.New()
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

func dictStream(o stream.Stream, n u.Node) stream.Stream {
	nextS := stream.New()
	for _, s := range o.O {
		// cartesian product
		partials := []map[string]any{
			{},
		}

		for _, c := range n.Children {

			innerS := stream.New()
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
