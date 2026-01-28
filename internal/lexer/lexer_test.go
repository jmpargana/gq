package lexer

import (
	"reflect"
	"testing"
)

func TestSimpleLexer(t *testing.T) {
	testCases := []struct {
		desc, input string
		tokens      []Token
	}{
		{
			desc:  "root",
			input: `.`,
			tokens: []Token{
				{Kind: DOT},
				{Kind: EOF},
			},
		},
		{
			desc:  "root iter",
			input: `.[]`,
			tokens: []Token{
				{Kind: DOT},
				{Kind: LBRACE},
				{Kind: RBRACE},
				{Kind: EOF},
			},
		},
		{
			desc:  "root index",
			input: `.[0]`,
			tokens: []Token{
				{Kind: DOT},
				{Kind: LBRACE},
				{Kind: NUMBER, Value: "0"},
				{Kind: RBRACE},
				{Kind: EOF},
			},
		},
		{
			desc:  "complex expression",
			input: `{b: [ ."a"[1].b.[1]] | .[0] }`,
			tokens: []Token{
				{Kind: LBRACKET},
				{Kind: IDENT, Value: "b"},
				{Kind: COLON},
				{Kind: LBRACE},
				{Kind: DOT},
				{Kind: STRING, Value: "a"},
				{Kind: LBRACE},
				{Kind: NUMBER, Value: "1"},
				{Kind: RBRACE},
				{Kind: DOT},
				{Kind: IDENT, Value: "b"},
				{Kind: DOT},
				{Kind: LBRACE},
				{Kind: NUMBER, Value: "1"},
				{Kind: RBRACE},
				{Kind: RBRACE},
				{Kind: PIPE},
				{Kind: DOT},
				{Kind: LBRACE},
				{Kind: NUMBER, Value: "0"},
				{Kind: RBRACE},
				{Kind: RBRACKET},
				{Kind: EOF},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := lex(tC.input)
			if !reflect.DeepEqual(got, tC.tokens) {
				t.Fatalf("expected: %v\ngot: %v\n", tC.tokens, got)
			}
		})
	}
}
