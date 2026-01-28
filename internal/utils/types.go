package utils

import (
	"strings"

	"github.com/jmpargana/gq/internal/gqjson"
)

type Kind int

const (
	INDEXSTART = iota
	INDEXEND
	FIELD
	IDX
	PIPE
	DICTSTART
	DICTEND
	ASSIGN
	COMMA
	ROOT
	ARRAY
)

type Cmd struct {
	Kind   Kind
	Fields []IdxField
	Ident  string
}

type IdxField struct {
	Name string
	Idx  int
	Kind Kind
}

type Node struct {
	Value    Cmd
	Children []Node
}

type Stream struct {
	O []any
}

func NewStream() Stream {
	return Stream{
		O: []any{},
	}
}

func NewSingleStream(obj any) Stream {
	return Stream{
		O: []any{obj},
	}
}

func (s *Stream) String() string {
	sb := strings.Builder{}
	for _, o := range s.O {
		sb.WriteString(gqjson.NewJSON(o).String())
	}
	return sb.String()
}
