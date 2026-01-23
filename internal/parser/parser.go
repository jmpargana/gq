package parser

import (
	u "github.com/jmpargana/gq/internal/utils"
)

type token u.Cmd

type Parser struct {
	ts  []token
	pos int
}

func NewParser(cs []u.Cmd) *Parser {
	ts := []token{}
	for _, c := range cs {
		ts = append(ts, token(c))
	}
	return &Parser{ts: ts, pos: 0}
}

func (p *Parser) peek() token {
	return p.ts[p.pos]
}

func (p *Parser) advance() token {
	t := p.ts[p.pos]
	p.pos++
	return t
}

func (p *Parser) match(kind u.Kind) bool {
	if p.pos >= len(p.ts) {
		return false
	}
	if p.peek().Kind == kind {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) expect(k u.Kind) token {
	if p.peek().Kind == k {
		return p.advance()
	}
	return token{}
}

func (p *Parser) ParseExpr() u.Node {
	term := p.parseTerm()

	for p.match(u.PIPE) {
		right := p.parseTerm()
		term = u.Node{Value: u.Cmd{Kind: u.PIPE}, Children: []u.Node{term, right}}
	}

	return term
}

func (p *Parser) parseTerm() u.Node {
	switch p.peek().Kind {
	case u.IDX:
		t := p.advance()
		return u.Node{Value: u.Cmd(t)}
	case u.INDEXSTART:
		t := p.advance()
		expr := p.ParseExpr()
		p.expect(u.INDEXEND)
		return u.Node{Value: u.Cmd(t), Children: []u.Node{expr}}
	case u.DICTSTART:
		return p.parseDict()
	}
	return u.Node{}
}

func (p *Parser) parseDict() u.Node {
	t := p.expect(u.DICTSTART)
	assignments := []u.Node{}

	assignments = append(assignments, p.parseAssignment())

	for p.match(u.COMMA) {
		assignments = append(assignments, p.parseAssignment())
	}

	p.expect(u.DICTEND)

	return u.Node{Value: u.Cmd(t), Children: assignments}
}

func (p *Parser) parseAssignment() u.Node {
	ident := p.expect(u.ASSIGN)
	return u.Node{Value: u.Cmd(ident), Children: []u.Node{p.ParseExpr()}}
}
