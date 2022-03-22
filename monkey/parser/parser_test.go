package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"strconv"
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
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

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}

}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	literal, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", il)
	}

	// Value は値そのもの
	if literal.Value != value {
		t.Errorf("ident.Value is not %d. got=%d", value, literal.Value)
	}

	// TokenLiteral は数値の場合は数値そのもの
	if literal.TokenLiteral() != strconv.FormatInt(value, 10) {
		t.Errorf("literal.TokenLiteral not %s. got=%s", strconv.FormatInt(value, 10), literal.TokenLiteral())
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program do not have enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not *ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestParsingInfixExpressionsMultipleInput1(t *testing.T) {
	infixTests := []struct {
		input            string
		prioriLeftValue  int64
		prioriOperator   string
		prioriRightValue int64
		operator         string
		RightValue       int64
	}{
		// expは
		// {RightValue} {operator} ( {prioriLeftValue} {prioriOperator} {prioriRightValue} )
		// という順で入っている
		{"5 + 4 * 2;", 4, "*", 2, "+", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program do not have enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not *ast.InfixExpression. got=%T", stmt.Expression)
		}

		// 優先度の最も低い RightValue を確認
		if !testIntegerLiteral(t, exp.Left, tt.RightValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		// 高優先度部分を分解
		exp2, ok := exp.Right.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp2 is not *ast.InfixExpression. got=%T", exp.Right)
		}

		if !testIntegerLiteral(t, exp2.Left, tt.prioriLeftValue) {
			return
		}

		if exp2.Operator != tt.prioriOperator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.prioriOperator, exp2.Operator)
		}

		if !testIntegerLiteral(t, exp2.Right, tt.prioriRightValue) {
			return
		}
	}
}

func TestParsingInfixExpressionsMultipleInput2(t *testing.T) {
	infixTests := []struct {
		input            string
		prioriLeftValue  int64
		prioriOperator   string
		prioriRightValue int64
		operator         string
		RightValue       int64
	}{
		// expは
		// ( {prioriLeftValue} {prioriOperator} {prioriRightValue} ) {operator} {RightValue}
		// という順で入っている
		{"5 * 4 + 2;", 5, "*", 4, "+", 2},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program do not have enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not *ast.InfixExpression. got=%T", stmt.Expression)
		}

		// 高優先度部分を分解
		exp2, ok := exp.Left.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp2 is not *ast.InfixExpression. got=%T", exp.Right)
		}

		if !testIntegerLiteral(t, exp2.Left, tt.prioriLeftValue) {
			return
		}

		if exp2.Operator != tt.prioriOperator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.prioriOperator, exp2.Operator)
		}

		if !testIntegerLiteral(t, exp2.Right, tt.prioriRightValue) {
			return
		}

		// 優先度の低い部分を処理
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		// 優先度の最も低い RightValue を確認
		if !testIntegerLiteral(t, exp.Right, tt.RightValue) {
			return
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a*b", "((-a) * b)"},
		{"-a", "(-a)"},
		{"-3278*vd+89*dv", "(((-3278) * vd) + (89 * dv))"},
		{"5<4!=3>4", "((5 < 4) != (3 > 4))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual_string := program.String()
		if actual_string != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, program.String())
		}
	}
}
