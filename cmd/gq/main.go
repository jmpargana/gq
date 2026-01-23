package main

import (
	"bufio"
	"os"

	"github.com/jmpargana/gq/internal/ast"
	json "github.com/jmpargana/gq/internal/gqjson"
	"github.com/jmpargana/gq/internal/lexer"
	"github.com/jmpargana/gq/internal/parser"
)

func main() {
	// FIXME: add error handlin
	r := bufio.NewReader(os.Stdin)
	tokens := lexer.Lex(os.Args[1])
	p := parser.NewParser(tokens)
	t := p.ParseExpr()
	obj := json.ParseObject(r)
	result := ast.Transform(obj, t)
	json.Print(result)
}
