package parser

import (
	"monkeylang/ast"
	"monkeylang/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
			let x = 5;
			let y = 10;
			let foobar = 838383;
			let barfoo 37;
			`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program produced %d statements instead of 3", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		currentStatement := program.Statements[i]
		if !testLetStatement(t, currentStatement, test.expectedIdentifer) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("Statement.TokenLiteral not let. got: %q", statement.TokenLiteral())
		return false
	}

	//This line asserts that the statement is a LetStatement (https://tour.golang.org/methods/15)
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement is not *ast.LetStatement. got: %T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("LetStatement name not: %s. got: %s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("LetStatement.Name.TokenLiteral() not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParseErrors(t *testing.T, parser *Parser) {
	errorMsgs := parser.Errors()

	numErrors := len(errorMsgs)
	if numErrors == 0 {
		return
	}

	t.Errorf("The parser encountered %d error(s)", numErrors)
	for _, errorMsg := range errorMsgs {
		t.Errorf("Parser Error: %s", errorMsg)
	}
	t.FailNow()
}
