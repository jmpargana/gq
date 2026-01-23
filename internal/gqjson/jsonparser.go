package gqjson

import (
	"bufio"
	"strconv"
	"strings"
)

func parseString(r *bufio.Reader) string {
	var sb strings.Builder
	var ignoreNext = false
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return sb.String()
		}
		if ch == '\\' {
			ignoreNext = true
		}
		if ignoreNext {
			ignoreNext = false
			continue
		}
		if ch == '"' {
			return sb.String()
		}
		sb.WriteRune(ch)
	}
}

func parseNumber(ch rune, r *bufio.Reader) any {
	b := []byte{byte(ch)}
	for {
		next, err := r.Peek(1)
		if err != nil || next[0] == ',' || next[0] == '}' || next[0] == ']' {
			break
		}
		ch, err := r.ReadByte()
		if err != nil {
			break
		}
		b = append(b, ch)
	}

	for _, c := range b {
		if c == '.' || c == 'e' || c == 'E' {
			f, _ := strconv.ParseFloat(string(b), 64)
			return f
		}
	}

	i, _ := strconv.ParseInt(string(b), 10, 64)
	return i
}

func parseBool(ch rune, r *bufio.Reader) bool {
	out := false
	n := 4
	if ch == 't' {
		n = 3
		out = true
	}
	for i := 0; i < n; i++ {
		_, _, err := r.ReadRune()
		if err != nil {
			panic(err)
		}
	}
	return out
}

func parseList(r *bufio.Reader) []any {
	out := []any{}
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			break
		}
		switch ch {
		case '{':
			out = append(out, ParseObject(r))
		case ',':
			continue
		case ']':
			return out
		case '[':
			out = append(out, parseList(r))
		case '"':
			out = append(out, parseString(r))
		case 't', 'f':
			out = append(out, parseBool(ch, r))
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			out = append(out, parseNumber(ch, r))
		}
	}
	return out
}

func ParseObject(r *bufio.Reader) any {
	out := map[string]any{}
	ident := ""
	key := true
	pendingKey := ""
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			break
		}
		switch ch {
		case '[':
			{
				if pendingKey != "" {
					out[pendingKey] = parseList(r)
				} else {
					return parseList(r)
				}
			}
		case '{':
			{
				if pendingKey != "" {
					out[pendingKey] = ParseObject(r)
				}
			}
		case '"':
			{
				ident = parseString(r)
				if key {
					pendingKey = ident
				} else {
					out[pendingKey] = ident
					pendingKey = ""
				}
			}
		case ':':
			{
				key = false
			}
		case ',':
			{
				key = true
			}
		case 't', 'f':
			{
				out[pendingKey] = parseBool(ch, r)
				pendingKey = ""
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			{
				out[pendingKey] = parseNumber(ch, r)
				pendingKey = ""
			}
		case ']':
			return out
		case '}':
			return out
		}
	}
	return out
}
