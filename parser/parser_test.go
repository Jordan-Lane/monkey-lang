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
			`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
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
		t.Errorf("statement.TokenLiteral not let. got: %q", statement.TokenLiteral())
		return false
	}

	//What does this line specifically do ???
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not *ast.statement. got: %T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("let statement name not: %s . got: %s", letStatement.Name.Value, name)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", letStatement.Name.TokenLiteral(), name)
		return false
	}

	return true
}
