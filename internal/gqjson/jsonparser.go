package gqjson

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseString(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	var ignoreNext = false
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return sb.String(), fmt.Errorf("failed parsing string after %s with: %v", sb.String(), err)
		}
		if ch == '\\' {
			ignoreNext = true
		}
		if ignoreNext {
			ignoreNext = false
			continue
		}
		if ch == '"' {
			return sb.String(), nil
		}
		sb.WriteRune(ch)
	}
}

func parseNumber(ch rune, r *bufio.Reader) (any, error) {
	b := []byte{byte(ch)}
	for {
		next, err := r.Peek(1)
		if err != nil || next[0] == ',' || next[0] == '}' || next[0] == ']' {
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("failed parsing number: %v", err)
			}
			break
		}
		ch, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed parsing number: %v", err)
			}
			break
		}
		b = append(b, ch)
	}

	for _, c := range b {
		if c == '.' || c == 'e' || c == 'E' {
			f, err := strconv.ParseFloat(string(b), 64)
			if err != nil {
				return nil, fmt.Errorf("failed parsing float: %v", err)
			}
			return f, nil
		}
	}

	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parsing int: %v", err)
	}
	return i, nil
}

func parseBool(ch rune, r *bufio.Reader) (bool, error) {
	out := false
	n := 4
	if ch == 't' {
		n = 3
		out = true
	}
	for i := 0; i < n; i++ {
		_, _, err := r.ReadRune()
		if err != nil {
			return false, fmt.Errorf("failed closing bool ident: %s", err)
		}
	}
	return out, nil
}

func parseList(r *bufio.Reader) ([]any, error) {
	out := []any{}
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed parsing list: %v", err)
			}
			break
		}
		switch ch {
		case '{':
			o, err := ParseObject(r)
			if err != nil {
				return nil, err
			}
			out = append(out, o)
		case ',':
			continue
		case ']':
			return out, nil
		case '[':
			l, err := parseList(r)
			if err != nil {
				return nil, err
			}
			out = append(out, l)
		case '"':
			s, err := parseString(r)
			if err != nil {
				return nil, err
			}
			out = append(out, s)
		case 't', 'f':
			b, err := parseBool(ch, r)
			if err != nil {
				return nil, err
			}
			out = append(out, b)
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			n, err := parseNumber(ch, r)
			if err != nil {
				return nil, err
			}
			out = append(out, n)
		}
	}
	return out, fmt.Errorf("unclosed array")
}

func ParseObject(r *bufio.Reader) (any, error) {
	out := map[string]any{}
	ident := ""
	key := true
	pendingKey := ""
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed parsing object: %v", err)
			}
			break
		}
		switch ch {
		case '[':
			if pendingKey != "" {
				l, err := parseList(r)
				if err != nil {
					return nil, err
				}
				out[pendingKey] = l
			} else {
				return parseList(r)
			}
		case '{':
			if pendingKey != "" {
				o, err := ParseObject(r)
				if err != nil {
					return nil, err
				}
				out[pendingKey] = o
			}
		case '"':
			s, err := parseString(r)
			if err != nil {
				return nil, err
			}
			ident = s
			if key {
				pendingKey = ident
			} else {
				out[pendingKey] = ident
				pendingKey = ""
			}
		case ':':
			key = false
		case ',':
			key = true
		case 't', 'f':
			b, err := parseBool(ch, r)
			if err != nil {
				return nil, err
			}
			out[pendingKey] = b
			pendingKey = ""
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			n, err := parseNumber(ch, r)
			if err != nil {
				return nil, err
			}
			out[pendingKey] = n
			pendingKey = ""
		case '}':
			return out, nil
		}
	}
	return out, fmt.Errorf("invalid object")
}
