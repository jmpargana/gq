package utils

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
