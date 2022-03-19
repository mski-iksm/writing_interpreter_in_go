package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota // _ を0にしてこのあとの定数に1から連番を振る
	LOWEST
	EQUALS      // ==
	LESSGREATER //> or <
	SUM         //+
	PRODUCT     //*
	PREFIX      //-X !X
	CALL        //function(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // 今のtoken
	peekToken token.Token // 次のtoken

	prefixParseFn map[token.TokenType]prefixParseFn // key が token.TokenType で value が prefixParseFn
	infixParseFn  map[token.TokenType]infixParseFn

	errors []string
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
func (p *Parser) ParseProgram() *ast.Program {
	// 空の Program structを新規作成
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{}
	stmt.Token = p.curToken // token.LET である

	// 変数名を期待
	if !p.expectPeek(token.IDENT) {
		return nil
		// 期待ハズレは全部nilを返す
	}
	// 変数名をセット
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// = を期待
	if !p.expectPeek(token.ASSIGN) {
		return nil
		// 期待ハズレは全部nilを返す
	}
	// = なら次へ行く

	// FIXME: セミコロンまで読み飛ばし
	// Valueに値が入っていない
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt

}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}
	stmt.Token = p.curToken

	p.NextToken()

	// FIXME: セミコロンまで読み飛ばし
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// FIXME: セミコロンまで読み飛ばし
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFn[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	}
	p.peekError(t) // エラーだよと主張する
	return false
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFn[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFn[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	val64, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = val64

	return lit
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// pの curToken の位置を動かすためにポインタで用意

	p.prefixParseFn = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)   // 変数名の関数をセットする
	p.registerPrefix(token.INT, p.parseIntegerLiteral) // intの関数をセットする

	// 2つトークンを読み込む
	// curToken と peekToken を読み込んでいる
	// 1回目だと peekがcurに入るだけなので2回呼ぶ必要がある
	p.NextToken()
	p.NextToken()

	return p
}
