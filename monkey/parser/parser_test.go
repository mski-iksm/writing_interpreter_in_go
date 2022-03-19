package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838883;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p) // parse checkはtestだけで実行するのでいいのか？

	if program == nil {
		// testのエラー
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements has wrond length. git=%d", len(program.Statements))
	}

	// listをリテラルで作成
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	// 3つのlet statementになっているかをチェック
	for i, tt := range tests {
		stmt := program.Statements[i]
		// identifier の名前(LetStatement.Name.Value)が合ってるかをチェック
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			// 空returnじゃなくてbreakでいいんじゃない？？
			break
		}
	}
}

// 小文字スタートなのでtestで直接は呼ばれない
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	// TokenLiteral() は *ast.Statement にしか実装されていないが
	// s.TokenLiteral() は自動的に (*s).TokenLiteral() と解釈されているために s からアクセスできる
	if s.TokenLiteral() != "let" {
		// そもそもletでないなら testのエラーをraiseする
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		// sが ast.LetStatement のインターフェースを持ってないとNG
		t.Errorf("s not ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s not=%s", name, letStmt.Name.Value)
		return false
	}

	// token.IDENT の literalは変数名（変数名の文字列そのもの）
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
		// %q: 対応する文字をシングルクォート'で囲んだ文字列
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p) // parse checkはtestだけで実行するのでいいのか？

	if program == nil {
		// testのエラー
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements has wrond length. got=%d", len(program.Statements))
	}

	// 3つの return statement になっているかをチェック
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		// stmt は interface
		// ast.ReturnStatement の型にしたものが returnStmt に入る

		// sが ast.ReturnStatement のインターフェースを持ってないとNG
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}

		// literal は "return" のはず
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}

	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program do not have enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	// インターフェイスが *ast.ExpressionStatement に実装されているかをチェック
	// されていれば stmt に *ast.ExpressionStatement の値をセットする
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
		// errofとの違いは？
	}

	// Value は変数名
	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not %s. got=%s", "foobar", ident.Value)
	}

	// TokenLiteral はidentの場合は変数名
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program do not have enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	// インターフェイスが *ast.ExpressionStatement に実装されているかをチェック
	// されていれば stmt に *ast.ExpressionStatement の値をセットする
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	// Value は値そのもの
	if literal.Value != 5 {
		t.Errorf("ident.Value is not %d. got=%d", 5, literal.Value)
	}

	// TokenLiteral は数値の場合は数値そのもの
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}
