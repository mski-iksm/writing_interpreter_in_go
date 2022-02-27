package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string // デバッグ用
}

// 文
type Statement interface {
	Node // Nodeを継承
	statementNode()
}

type Expression interface {
	Node // Nodeを継承
	expressionNode()
}

// Root Node
type Program struct {
	Statements []Statement // 文が格納されるリスト
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// let statement
type LetStatement struct {
	Token token.Token // token.LET というtokenを格納する
	Name  *Identifier // 変数名; なぜポインタ???
	Value Expression  // 格納する式
}

// Statementのインターフェースを揃える
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token // token.RETURN というtokenを格納する
	ReturnValue Expression  // 格納する式
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type Identifier struct {
	Token token.Token // token.IDENT というtokenを保持Token
	Value string
}

// Expressionのインターフェースを揃える
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
