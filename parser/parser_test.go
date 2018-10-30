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
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program produced %d statements instead of 3", len(program.Statements))
	}

	letStatementTests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range letStatementTests {
		currentStatement := program.Statements[i]
		if !testLetStatement(t, currentStatement, test.expectedIdentifer) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("Statement.TokenLiteral not let. Got: %q", statement.TokenLiteral())
		return false
	}

	//This line asserts that the statement is a LetStatement (https://tour.golang.org/methods/15)
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement is not *ast.LetStatement. Got: %T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("LetStatement.Name.Value not: %s. Got: %s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("LetStatement.Name.TokenLiteral() not '%s'. Got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	// Later, test the actual expression

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
			return 5;
			return x;
			return add(x+y);
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

	for _, currentStatement := range program.Statements {
		if !testReturnStatement(t, currentStatement) {
			return
		}
	}
}

func testReturnStatement(t *testing.T, statement ast.Statement) bool {
	// Check that the statement token is a return
	if statement.TokenLiteral() != "return" {
		t.Errorf("Statement.TokenLiteral not return. Got %q", statement.TokenLiteral())
		return false
	}

	// Assert that the statement is a return statement
	_, ok := statement.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("Statement is not *ast.ReturnStatement. Got %T", statement)
		return false
	}

	// Later, test the actual expression
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
