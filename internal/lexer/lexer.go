package lexer

import (
	"bufio"
	"strings"
	"unicode"
)

type TokenKind int

const (
	LBRACKET TokenKind = iota
	RBRACKET
	LBRACE
	RBRACE
	DOT
	PIPE
	COMMA
	COLON
	IDENT
	NUMBER
	STRING
	EOF
	ILLEGAL
)

type Token struct {
	Kind  TokenKind
	Value string
}

type Lexer struct {
	r   *bufio.Reader
	ch  rune
	eof bool
}

func Lex(s string) []Token {
	return lex(s)
}

func lex(s string) []Token {
	l := newLexer(s)
	var tokens []Token
	for {
		tok := l.nextToken()
		tokens = append(tokens, tok)
		if tok.Kind == EOF {
			break
		}
	}
	return tokens
}

func newLexer(s string) *Lexer {
	r := bufio.NewReader(strings.NewReader(s))
	ch, _, err := r.ReadRune()
	if err != nil {
		return &Lexer{
			r:   r,
			ch:  0,
			eof: true,
		}
	}
	return &Lexer{
		r:   r,
		ch:  ch,
		eof: false,
	}
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case 0:
		l.eof = true
		return Token{Kind: EOF}
	case '.':
		l.read()
		return Token{Kind: DOT}
	case ',':
		l.read()
		return Token{Kind: COMMA}
	case ':':
		l.read()
		return Token{Kind: COLON}
	case '|':
		l.read()
		return Token{Kind: PIPE}
	case '{':
		l.read()
		return Token{Kind: LBRACKET}
	case '}':
		l.read()
		return Token{Kind: RBRACKET}
	case '[':
		l.read()
		return Token{Kind: LBRACE}
	case ']':
		l.read()
		return Token{Kind: RBRACE}
	default:
		if isDigit(l.ch) {
			return l.readNumber()
		}
		if isIdentStart(l.ch) {
			return l.readIdent()
		}
		if l.ch == '"' {
			return l.readString()
		}
		illegal := l.ch
		l.read()
		return Token{Kind: ILLEGAL, Value: string(illegal)}
	}
}

func (l *Lexer) readNumber() Token {
	var b strings.Builder
	for isDigit(l.ch) {
		b.WriteRune(l.ch)
		l.read()
	}
	return Token{Kind: NUMBER, Value: b.String()}
}

func (l *Lexer) readString() Token {
	l.read() // skip "
	var b strings.Builder
	for l.ch != '"' && l.ch != 0 {
		b.WriteRune(l.ch)
		l.read()
	}
	l.read() // skip "
	return Token{Kind: STRING, Value: b.String()}
}

func (l *Lexer) readIdent() Token {
	var b strings.Builder
	for isIdentChar(l.ch) {
		b.WriteRune(l.ch)
		l.read()
	}
	return Token{Kind: IDENT, Value: b.String()}
}

func (l *Lexer) read() {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		l.ch = 0
		return
	}
	l.ch = ch
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' {
		l.read()
	}
}

// TODO: move to utils
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentChar(r rune) bool {
	return isIdentStart(r) || isDigit(r)
}
