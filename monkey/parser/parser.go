package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // 今のtoken
	peekToken token.Token // 次のtoken
	errors    []string
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
	// 空の Ptogram structを新規作成
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
		return nil
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

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// pの curToken の位置を動かすためにポインタで用意

	p.NextToken()
	p.NextToken()

	return p
}
