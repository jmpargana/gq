package main

import "fmt"

const ident = 2

var reset = "\x1b[0m"

var levelColor = map[int]string{
	0: "\x1b[31",
	1: "\x1b[32",
	2: "\x1b[33",
	3: "\x1b[34",
	4: "\x1b[35",
}

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

func print(s any) {
	switch s := s.(type) {
	case map[string]any:
		fmt.Println(printObj(s, 0))
	case []any:
		fmt.Println(printList(s, 0))
	}
}
