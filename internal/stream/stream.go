package stream

import (
	"strings"

	"github.com/jmpargana/gq/internal/gqjson"
)

type Stream struct {
	O []any
}

func New() Stream {
	return Stream{
		O: []any{},
	}
}

func NewS(obj any) Stream {
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
