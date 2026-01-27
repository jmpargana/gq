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

type Stream struct {
	O []any
}

func NewStream() Stream {
	return Stream{
		O: []any{},
	}
}

func (c Cmd) IsEqual(o Cmd) bool {
	if c.Kind != o.Kind {
		return false
	}
	if c.Kind == IDX {
		if len(c.Fields) != len(o.Fields) {
			return false
		}

		for i := range c.Fields {
			if c.Fields[i] != o.Fields[i] {
				return false
			}
		}
	}
	return true
}

func (n Node) IsEqual(o Node) bool {
	if !n.Value.IsEqual(o.Value) || len(n.Children) != len(o.Children) {
		return false
	}
	if len(n.Children) > 0 {
		for i := range n.Children {
			if !n.Children[i].IsEqual(o.Children[i]) {
				return false
			}
		}
	}
	return true
}
