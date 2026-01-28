package gqjson

import (
	"fmt"
	"strings"
)

type JSON struct {
	O any
}

func NewJSON(o any) *JSON {
	return &JSON{O: o}
}

func (j *JSON) String() string {
	sb := strings.Builder{}
	switch s := j.O.(type) {
	case map[string]any:
		sb.WriteString(printObj(s, 0))
	case []any:
		sb.WriteString(printList(s, 0))
	case int64:
		fmt.Fprintf(&sb, "%d", s)
	case bool:
		fmt.Fprintf(&sb, "%t", s)
	case float64:
		fmt.Fprintf(&sb, "%0.2f", s)
	case string:
		fmt.Fprintf(&sb, "\"%s\"", s)
	}
	sb.WriteRune('\n')
	return sb.String()
}

const ident = 2

// FIXME: fix performance issue with string append
func printList(l []any, level int) string {
	s := "[\n"
	for i, it := range l {
		for range (level + 1) * ident {
			s += " "
		}
		switch it := it.(type) {
		case string:
			s += fmt.Sprintf("\"%s\"", it)
		case int, int16, int32, int64, int8:
			s += fmt.Sprintf("%d", it)
		case float64, float32:
			s += fmt.Sprintf("%f", it)
		case bool:
			s += fmt.Sprintf("%t", it)
		case map[string]any:
			s += printObj(it, level+1)
		case []any:
			s += printList(it, level+1)
		}
		if i < len(l)-1 {
			s += ","
		}
		s += "\n"
	}
	for range level * ident {
		s += " "
	}
	s += "]"
	return s
}

func printObj(obj map[string]any, level int) string {
	s := ""
	s += "{\n"

	i := 0
	n := len(obj)

	for k, v := range obj {
		for range (level + 1) * ident {
			s += " "
		}
		s += fmt.Sprintf("\"%s\": ", k)
		switch v := v.(type) {
		case string:
			s += fmt.Sprintf("\"%s\"", v)
		case int, int16, int32, int64, int8:
			s += fmt.Sprintf("%d", v)
		case float32, float64:
			s += fmt.Sprintf("%.2f", v)
		case bool:
			s += fmt.Sprintf("%t", v)
		case map[string]any:
			s += printObj(v, level+1)
		case []any:
			s += printList(v, level+1)
		}
		i++
		if i < n {
			s += ","
		}
		s += "\n"
	}

	for range level * ident {
		s += " "
	}

	s += "}"
	return s
}
