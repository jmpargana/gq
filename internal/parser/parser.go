package parser

import (
	"strconv"

	"github.com/jmpargana/gq/internal/lexer"
	u "github.com/jmpargana/gq/internal/utils"
)

type Parser struct {
	ts  []lexer.Token
	pos int
}

func NewParser(cs []lexer.Token) *Parser {
	return &Parser{ts: cs, pos: 0}
}

func (p *Parser) peek() lexer.Token {
	return p.ts[p.pos]
}

func (p *Parser) advance() lexer.Token {
	t := p.ts[p.pos]
	p.pos++
	return t
}

func (p *Parser) match(kind lexer.TokenKind) bool {
	if p.pos >= len(p.ts) {
		return false
	}
	if p.peek().Kind == kind {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) expect(k lexer.TokenKind) lexer.Token {
	if p.peek().Kind == k {
		return p.advance()
	}
	return lexer.Token{}
}

func (p *Parser) ParseExpr() u.Node {
	term := p.parseTerm()

	for p.match(lexer.PIPE) {
		right := p.parseTerm()
		term = u.Node{Value: u.Cmd{Kind: u.PIPE}, Children: []u.Node{term, right}}
	}

	return term
}

func (p *Parser) parseTerm() u.Node {
	switch p.peek().Kind {
	case lexer.DOT:
		return p.parseIndex()
	case lexer.LBRACE:
		p.advance()
		expr := p.ParseExpr()
		p.expect(lexer.RBRACE)
		return u.Node{Value: u.Cmd{Kind: u.INDEXSTART}, Children: []u.Node{expr}}
	case lexer.LBRACKET:
		return p.parseDict()
	// FIXME:
	default:
		return u.Node{}
	}
}

func (p *Parser) parseIndex() u.Node {
	idxs := []u.IdxField{}

	p.expect(lexer.DOT)
	for {
		tok := p.peek().Kind
		if !isValidIndexStarter(tok) {
			if len(idxs) == 0 {
				idxs = append(idxs, u.IdxField{Kind: u.ROOT})
			}
			return u.Node{Value: u.Cmd{Kind: u.IDX, Fields: idxs}}
		}

		switch tok {
		case lexer.DOT:
			p.advance()
		case lexer.IDENT, lexer.STRING:
			t := p.advance()
			idxs = append(idxs, u.IdxField{Kind: u.FIELD, Name: t.Value})
		case lexer.LBRACE:
			p.advance()
			switch p.peek().Kind {
			case lexer.RBRACE:
				p.advance()
				idxs = append(idxs, u.IdxField{Kind: u.ARRAY})
			case lexer.NUMBER:
				t := p.advance()
				n, err := strconv.Atoi(t.Value)
				if err != nil {
					panic(err)
				}
				if !p.match(lexer.RBRACE) {
					panic("need to close brace")
				}
				idxs = append(idxs, u.IdxField{Kind: u.IDX, Idx: n})
			case lexer.IDENT, lexer.STRING:
				t := p.advance()
				if !p.match(lexer.RBRACE) {
					panic("need to close brace")
				}
				idxs = append(idxs, u.IdxField{Kind: u.FIELD, Name: t.Value})
			// FIXME: break if not matching
			default:
				continue
			}
		// FIXME: break if not matching
		default:
			continue
		}
	}

}

func (p *Parser) parseDict() u.Node {
	p.expect(lexer.LBRACKET)
	assignments := []u.Node{}

	assignments = append(assignments, p.parseAssignment())

	for p.match(lexer.COMMA) {
		assignments = append(assignments, p.parseAssignment())
	}

	p.expect(lexer.RBRACKET)

	return u.Node{Value: u.Cmd{Kind: u.DICTSTART}, Children: assignments}
}

func (p *Parser) parseAssignment() u.Node {
	ident := p.expect(lexer.IDENT)
	p.expect(lexer.COLON)
	return u.Node{Value: u.Cmd{Kind: u.ASSIGN, Ident: ident.Value}, Children: []u.Node{p.ParseExpr()}}
}

func isValidIndexStarter(t lexer.TokenKind) bool {
	return t == lexer.IDENT || t == lexer.STRING || t == lexer.LBRACE || t == lexer.DOT
}
