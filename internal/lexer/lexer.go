package lexer

import (
	"bufio"
	"strconv"
	"strings"

	u "github.com/jmpargana/gq/internal/utils"
)

func createIdxField(s string) u.IdxField {
	n, err := strconv.Atoi(s)
	if err != nil {
		if s == "" {
			return u.IdxField{Kind: u.ROOT}
		}
		return u.IdxField{Kind: u.FIELD, Name: s}
	}
	return u.IdxField{Kind: u.IDX, Idx: n}
}

func Lex(program string) []u.Cmd {
	// TODO: missing all array index => []
	pgr := []u.Cmd{}
	fields := []u.IdxField{}
	r := bufio.NewReader(strings.NewReader(program))
	fieldName := ""
	indexing := false
	wrapping := 0
	mapping := 0
	prev := ' '
	for {
		rn, _, err := r.ReadRune()
		if err != nil {
			break
		}
		switch rn {
		case '|':
			if indexing {
				fields = append(fields, createIdxField(fieldName))
				pgr = append(pgr, u.Cmd{Kind: u.IDX, Fields: fields})
				fields = []u.IdxField{}
				fieldName = ""
				indexing = false
				pgr = append(pgr, u.Cmd{Kind: u.PIPE})
			} else {
				fields = []u.IdxField{}
				fieldName = ""
				indexing = false
				pgr = append(pgr, u.Cmd{Kind: u.PIPE})
			}
		case '"', ' ':
		case ']':
			// Check conflict with }
			if wrapping > 0 {
				if prev != ']' && prev != '}' {
					fields = append(fields, createIdxField(fieldName))
					pgr = append(pgr, u.Cmd{Kind: u.IDX, Fields: fields})
				}
				fieldName = ""
				fields = []u.IdxField{}
				pgr = append(pgr, u.Cmd{Kind: u.INDEXEND})
				wrapping--
				indexing = false
				prev = ']'
			}
		case '[':
			if !indexing {
				pgr = append(pgr, u.Cmd{Kind: u.INDEXSTART})
				wrapping++
			}
			if fieldName != "" {
				fields = append(fields, createIdxField(fieldName))
				fieldName = ""
			}
		case '.':
			indexing = true
			if fieldName != "" {
				fields = append(fields, createIdxField(fieldName))
				fieldName = ""
			}
		case '{':
			mapping++
			pgr = append(pgr, u.Cmd{Kind: u.DICTSTART})
		case ':':
			if fieldName != "" {
				pgr = append(pgr, u.Cmd{Kind: u.ASSIGN, Ident: fieldName})
				fieldName = ""
			}
		case ',':
			if prev != ']' && prev != '}' {
				fields = append(fields, createIdxField(fieldName))
				pgr = append(pgr, u.Cmd{Kind: u.IDX, Fields: fields})
				prev = ' '
			}
			fieldName = ""
			fields = []u.IdxField{}
			pgr = append(pgr, u.Cmd{Kind: u.COMMA})
		case '}':
			if mapping > 0 {
				if prev != ']' && prev != '}' {
					fields = append(fields, createIdxField(fieldName))
					pgr = append(pgr, u.Cmd{Kind: u.IDX, Fields: fields})
				}
				fieldName = ""
				fields = []u.IdxField{}
				pgr = append(pgr, u.Cmd{Kind: u.DICTEND})
				mapping--
				indexing = false
				prev = '}'
			}
		default:
			// TODO: change to string builder
			fieldName += string(rn)
			prev = ' '
		}
	}

	// TODO: fix last value pop
	if len(pgr) == 0 || (pgr[len(pgr)-1].Kind != u.INDEXEND && pgr[len(pgr)-1].Kind != u.DICTEND) {
		fields = append(fields, createIdxField(fieldName))
		pgr = append(pgr, u.Cmd{Kind: u.IDX, Fields: fields})
	}

	return pgr
}
