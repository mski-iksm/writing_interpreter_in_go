package ast

import (
	"bytes"
	"monkey/token"
	"strconv"
)

type Node interface {
	TokenLiteral() string // デバッグ用
	String() string       // デバッグでASTノードを表示用
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
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// let statement
type LetStatement struct {
	Token token.Token // token.LET というtokenを格納する
	Name  *Identifier // 変数名; なぜポインタ???
	Value Expression  // 格納する式
}

// LetStatement のインターフェースを揃える
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// return statement
type ReturnStatement struct {
	Token       token.Token // token.RETURN というtokenを格納する
	ReturnValue Expression  // 格納する式
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// 式 statement
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression  // 残りの式
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// ===================

// let で使う変数名
type Identifier struct {
	Token token.Token // token.IDENT というtokenを保持Token
	Value string
}

// Expressionのインターフェースを揃える
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

// 整数
type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
	// Token.Literal は もともと全部string
}
func (il *IntegerLiteral) String() string {
	return strconv.FormatInt(il.Value, 10)
}
