package ast

import (
	"fmt"
	"strings"

	u "github.com/jmpargana/gq/internal/utils"
)

func PrintAST(n u.Node, ident int) string {
	s := strings.Builder{}

	for i := 0; i < ident; i++ {
		s.WriteRune(' ')
	}

	s.WriteString(printCmd(n.Value, ident))
	// s.WriteRune('\n')

	for _, c := range n.Children {
		s.WriteString(PrintAST(c, ident+2))
	}

	return s.String()
}

func printCmd(c u.Cmd, ident int) string {
	s := strings.Builder{}

	switch c.Kind {
	case u.ASSIGN:
		fmt.Fprintf(&s, "ASSIGN: %s", c.Ident)
	case u.ARRAY:
		fmt.Fprintf(&s, "ARRAY:")
	case u.DICTSTART:
		fmt.Fprintf(&s, "DICT:")
	case u.PIPE:
		fmt.Fprintf(&s, "PIPE:")
	case u.IDX:
		fmt.Fprintf(&s, "IDX:")
		for _, f := range c.Fields {
			s.WriteRune('\n')
			for i := 0; i < ident+2; i++ {
				s.WriteRune(' ')
			}
			switch f.Kind {
			case u.IDX:
				fmt.Fprintf(&s, "LIST: %d", f.Idx)
			case u.ROOT:
				fmt.Fprintf(&s, "ROOT")
			case u.ARRAY:
				fmt.Fprintf(&s, "ITER")
			case u.FIELD:
				fmt.Fprintf(&s, "FIELD: %s", f.Name)
			}
		}
	}

	s.WriteRune('\n')

	return s.String()
}
