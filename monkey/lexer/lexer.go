package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int  // 現在の文字chの位置
	readPosition int  // これから読み込む文字の位置
	ch           byte // 現在検査中の文字
}

func (l *Lexer) readChar() {
	// ポインタレシーバを使うことでlの中身を変更することができる
	// 普通のレシーバだとlのコピーを触ることになるので変更が反映されない
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCIIコードのNULLに対応
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peakChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // ASCIIコードのNULLに対応
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	// 文字列を読み込んでTokenにして返す

	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		// 次の文字を先読みして、`==`となっているならRQにする
		if l.peakChar() == '=' {
			l.readChar()
			tok = newTokenFromString(token.EQ, "==")
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		// 次の文字を先読みして、`!=`となっているならRQにする
		if l.peakChar() == '=' {
			l.readChar()
			tok = newTokenFromString(token.NOT_EQ, "!=")
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()          // 文字列の塊を取得
			tok.Type = token.LookupIdent(tok.Literal) // keywords（予約語かを判定）
			return tok
			// ここは1文字進める必要がないための措置
			// readIdentifierの最後でreadChar()しているからだけどあんまりよくない気がする
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar() // 1文字すすめる
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { // l.chはループのたびに値が変わり、再評価される
		l.readChar() // readCharを使うことで終わりになるまで次々と文字を読んでいく
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) { // l.chはループのたびに値が変わり、再評価される
		l.readChar() // readCharを使うことで終わりになるまで次々と文字を読んでいく
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar() // 空白をスキップする
	}
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	// &でポインタを返すようにする
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newTokenFromString(tokenType token.TokenType, str string) token.Token {
	return token.Token{Type: tokenType, Literal: str}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
