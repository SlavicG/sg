package lexer

import "sg_interpreter/src/sg/token"

type Lexer struct {
	input string
	pos   int  //current pos
	nxt   int  //next pos
	ch    byte //cur char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	var tok token.Token

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: lit}
		} else {
			tok = newToken(token.SET, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.STAR, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '(':
		tok = newToken(token.LP, l.ch)
	case ')':
		tok = newToken(token.RP, l.ch)
	case '{':
		tok = newToken(token.LB, l.ch)
	case '}':
		tok = newToken(token.RB, l.ch)
	case '!':
		if l.peek() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: lit}
		} else {
			tok = newToken(token.EXC, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOL, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '[':
		tok = newToken(token.LBP, l.ch)
	case ']':
		tok = newToken(token.RBP, l.ch)
	case ':':
		tok = newToken(token.COL, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = token.FindIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNum()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.nxt >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nxt]
	}
	l.pos = l.nxt
	l.nxt++
}

func (l *Lexer) readString() string {
	position := l.pos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.pos]
}

func (l *Lexer) peek() byte {
	if l.nxt >= len(l.input) {
		return 0
	} else {
		return l.input[l.nxt]
	}
}

func (l *Lexer) readNum() string {
	position := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readIdent() string {
	position := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func isLetter(ch byte) bool {

	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch == '_')
}

func isDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9')
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
