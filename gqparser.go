package main

import (
	"bufio"
	"strconv"
	"strings"
)

type Kind int

const (
	mapIndexStart = iota
	mapIndexEnd
	// chain
	// array
	// dict
	field
	idx
	pipe
)

type cmd struct {
	kind   Kind
	fields []idxField
}

type idxField struct {
	name string
	idx  int
	kind Kind
}

type node struct {
	value    cmd
	children []node
}

func (c cmd) isEqual(o cmd) bool {
	if c.kind != o.kind {
		return false
	}
	if c.kind == idx {
		if len(c.fields) != len(o.fields) {
			return false
		}

		for i := range c.fields {
			if c.fields[i] != o.fields[i] {
				return false
			}
		}
	}
	return true
}

func createIdxField(s string) idxField {
	n, err := strconv.Atoi(s)
	if err != nil {
		return idxField{kind: field, name: s}
	}
	return idxField{kind: idx, idx: n}
}

func parseGQ(program string) []cmd {
	// TODO: missing all array index => []
	pgr := []cmd{}
	fields := []idxField{}
	r := bufio.NewReader(strings.NewReader(program))
	fieldName := ""
	indexing := false
	wrapping := false
	for {
		rn, _, err := r.ReadRune()
		if err != nil {
			break
		}
		switch rn {
		case '|':
			// TODO: handle existing
			fields = append(fields, createIdxField(fieldName))
			pgr = append(pgr, cmd{kind: idx, fields: fields})
			fields = []idxField{}
			fieldName = ""
			indexing = false
			pgr = append(pgr, cmd{kind: pipe})
		case '"', ' ':
		case ']':
			if wrapping {
				fields = append(fields, createIdxField(fieldName))
				pgr = append(pgr, cmd{kind: idx, fields: fields})
				fieldName = ""
				fields = []idxField{}
				pgr = append(pgr, cmd{kind: mapIndexEnd})
				wrapping = false
			}
		case '[':
			if !indexing {
				pgr = append(pgr, cmd{kind: mapIndexStart})
				wrapping = true
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
		default:
			// TODO: change to string builder
			fieldName += string(rn)
		}
	}

	// TODO: fix last value pop
	if len(pgr) == 0 || pgr[len(pgr)-1].kind != mapIndexEnd {
		fields = append(fields, createIdxField(fieldName))
		pgr = append(pgr, cmd{kind: idx, fields: fields})
	}

	return pgr
}

func transform(a any, program []cmd) any {
	prev := a
	for _, c := range program {
		if c.kind == idx {
			for _, f := range c.fields {
				if f.kind == idx {
					l := prev.([]any)
					prev = l[f.idx]
				}
				if f.kind == field {
					m := prev.(map[string]any)
					prev = m[f.name]
				}
			}
		}
	}
	return prev
}

func buildTree(pgr []cmd) node {
	if len(pgr) == 1 {
		return node{value: pgr[0]}
	}

	c := pgr[0]
	n := node{}

	switch c.kind {
	case mapIndexStart:
		for i := len(pgr) - 1; i > 0; i-- {
			if pgr[i].kind == mapIndexEnd {
				pgr = append(pgr[:i], pgr[i+1:]...)
			}
		}
		n.value = c
		n.children = append(n.children, buildTree(pgr[1:]))
	case idx:
		if len(pgr) >= 3 && pgr[1].kind == pipe {
			n.value = pgr[1]
			n.children = append(n.children, node{value: pgr[0]}, buildTree(pgr[2:]))
		} else {
			n.value = pgr[0]
		}
		// mapIndexEnd and pipe should be unreachable
	}
	return n
}

// func executeTree(a any, pgr tree) any { return nil }
