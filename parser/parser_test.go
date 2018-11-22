package parser

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/lexer"
	"strconv"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
		}

		testLetStatement(t, program.Statements[0], test.input, test.expectedIdentifier, test.expectedValue)
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, input string, expectedName string, value interface{}) bool {
	//This line asserts that the statement is a LetStatement (https://tour.golang.org/methods/15)
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement type is incorrect. Expected: *ast.LetStatement. Got: %T", statement)
		return false
	}

	if statement.TokenLiteral() != "let" {
		t.Errorf("Statement.TokenLiteral not let. Got: %q", statement.TokenLiteral())
		return false
	}

	if letStatement.Name.Value != expectedName {
		t.Errorf("LetStatement.Name.Value not: %s. Got: %s", expectedName, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != expectedName {
		t.Errorf("LetStatement.Name.TokenLiteral() not '%s'. Got=%s", expectedName, letStatement.Name.TokenLiteral())
		return false
	}

	if !testLiteralExpression(t, letStatement.Value, value) {
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return x;", "x"},
		{"return true;", true},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
		}

		testReturnStatement(t, program.Statements[0], test.expectedValue)
	}
}

func testReturnStatement(t *testing.T, expression ast.Statement, expectedValue interface{}) bool {

	returnStatement, ok := expression.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ReturnStatement. Got: %T", expression)
		return false
	}

	// Check that the statement token is a return
	expectedToken := "return"
	if returnStatement.TokenLiteral() != expectedToken {
		t.Errorf("Statement.TokenLiteral is incorrect. Expected: %q. Got %q", expectedToken, returnStatement.TokenLiteral())
		return false
	}

	if !testLiteralExpression(t, returnStatement.ReturnValue, expectedValue) {
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", statement)
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Statement.Expression is incorrect. Expected: *ast.Identifier. Got %T", statement.Expression)
	}
	expectedIdent := "foobar"
	if ident.Value != expectedIdent {
		t.Errorf("ident.Value is incorrect. Expected: %q. Got %q", expectedIdent, ident)
	}
	if ident.TokenLiteral() != expectedIdent {
		t.Errorf("ident.TokenLiteral is incorrect. Expected: %q. Got %q", expectedIdent, ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", statement)
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Statement.Expression is incorrect. Expected: *ast.IntegerLiteral. Got %T", statement.Expression)
	}

	expectedLiteral := "5"
	expectedInt, _ := strconv.ParseInt(expectedLiteral, 0, 64)
	if literal.Value != expectedInt {
		t.Errorf("integer.Value is incorrect. Expected: %d. Got %q", expectedInt, literal)
	}
	if literal.TokenLiteral() != expectedLiteral {
		t.Errorf("integer.TokenLiteral is incorrect. Expected: %q. Got %q", expectedLiteral, literal.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", statement)
	}

	stringLiteral, ok := statement.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Statement.Expression is incorrect. Expected: *ast.StringLiteral. Got: %T", statement.Expression)
	}

	if stringLiteral.Value != "hello world" {
		t.Fatalf("stringLiteral.Value is incorrect. Expected: hello world. Got: %s", stringLiteral.Value)
	}

}

func TestBooleanExpression(t *testing.T) {
	input := "true"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instaed of 1", len(program.Statements))
	}

	expression, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	booleanLiteral, ok := expression.Expression.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("Statement.Expression type is incorrect. Expected: *ast.BooleanLiteral. Got %T", expression.Expression)
	}

	expectedLiteral := "true"
	if booleanLiteral.TokenLiteral() != expectedLiteral {
		t.Errorf("BooleanLiteral.TokenLiteral is incorrect. Expected: %q. Got %q", expectedLiteral, booleanLiteral.TokenLiteral())
	}

	expectedBool, _ := strconv.ParseBool(input)
	if booleanLiteral.Value != expectedBool {
		t.Errorf("BooleanLiteral.Value is incorrect. Expected: %t. Got: %t", expectedBool, booleanLiteral.Value)
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-10", "-", 10},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, test := range prefixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement type is incorrect. Expected: *ast.Expression. Got: %T", statement)
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Statement.Expression is incorrect. Expected: *ast.PrefixExpression. Got %T", statement.Expression)
		}
		if expression.Operator != test.operator {
			t.Fatalf("Operator is incorrect. Expected: %s. Got %s", test.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, test.value) {
			t.Fatalf("Integer value is incorrect. Expected: %d. Got %d", test.value, expression.Right)
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"false == false", false, "==", false},
		{"true != false", true, "!=", false},
		{"false != true", false, "!=", true},
	}

	for _, test := range infixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement type is incorrect. Expected: *ast.Expression. Got: %T", statement)
		}

		if !testInfixExpression(t, statement.Expression, test.leftValue, test.operator, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b;",
			"((-a) * b)",
		},
		{
			"!-a;",
			"(!(-a))",
		},
		{
			"a + b + c;",
			"((a + b) + c)",
		},
		{
			"a + b - c;",
			"((a + b) - c)",
		},
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		}, {
			"3 + 4; -5 * 5;",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
		}, {
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"(1 + 2) + 3",
			"((1 + 2) + 3)",
		},
		{
			"1 + (2 + 3)",
			"(1 + (2 + 3))",
		},
		{
			"1 + (2 + 3)",
			"(1 + (2 + 3))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)
		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected=%q, got=%q", test.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x; }"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Statement.Expression type is incorrect. Expected: *ast.IfExpression. Got: %T", statement)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Fatalf("IfExpression.Consquence produced %d statements instead of 1", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statement[0] type is incorrect. Expected *ast.ExpressionStatement. Got %T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("IfExpression.Alternative is incorrect. Expected: nil. Got: %T", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x; } else { y; }"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Statement.Expression type is incorrect. Expected: *ast.IfExpression. Got: %T", statement)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Fatalf("IfExpression.Consquence produced %d statements instead of 1", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statement[0] type is incorrect. Expected *ast.ExpressionStatement. Got %T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("IfExpression.Alternative type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", expression.Alternative)
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "fn( x, y ){ x + y; }"
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Statement.Expression type is incorrect. Expected: *ast.FunctionLitearl. Got: %T", statement)
	}

	if len(expression.Parameters) != 2 {
		t.Fatalf("Number of parameters is incorrect. Expected: 2. Got: %d", len(expression.Parameters))
	}

	if !testLiteralExpression(t, expression.Parameters[0], "x") {
		return
	}
	if !testLiteralExpression(t, expression.Parameters[1], "y") {
		return
	}

	if len(expression.Body.Statements) != 1 {
		t.Fatalf("Expression.Body number of statements is incorrect. Expected: 1. Got: %d", len(expression.Body.Statements))
	}

	bodyStatement, ok := expression.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Body.Statement[0] type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", expression.Body.Statements[0])
	}

	if !testInfixExpression(t, bodyStatement.Expression, "x", "+", "y") {
		return
	}
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input              string
		expectedParameters []string
	}{
		{input: "fn(){};", expectedParameters: []string{}},
		{input: "fn(x){};", expectedParameters: []string{"x"}},
		{input: "fn(x, y){};", expectedParameters: []string{"x", "y"}},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(test.expectedParameters) {
			t.Fatalf("Number of parameters is incorrect. Expected: %d. Got %d", len(test.expectedParameters), len(function.Parameters))
		}

		for i, ident := range test.expectedParameters {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestFunctionCall(t *testing.T) {
	input := " add( 2 + 2, 5 * 5, 7)"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement type is incorrect. Expected: *ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	functionCall, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Statement.Expression type is incorrect. Expected: *ast.CallFunction. Got %T", statement.Expression)
	}

	if !testIdentifier(t, functionCall.Function, "add") {
		return
	}

	if len(functionCall.Arguments) != 3 {
		t.Fatalf("Number of FunctionCall.Arguments is incorrect. Expected: 3. Got: %d", len(functionCall.Arguments))
	}

	testInfixExpression(t, functionCall.Arguments[0], 2, "+", 2)
	testInfixExpression(t, functionCall.Arguments[1], 5, "*", 5)
	testLiteralExpression(t, functionCall.Arguments[2], 7)

}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {

	// This is a type switch (https://tour.golang.org/methods/16)
	switch value := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(value))
	case int64:
		return testIntegerLiteral(t, expression, value)
	case string:
		return testIdentifier(t, expression, value)
	case bool:
		return testBooleanLiteral(t, expression, value)
	}

	t.Errorf("Expression literal not supported: %T", expression)
	return false
}

func testIntegerLiteral(t *testing.T, integerLiteral ast.Expression, expectedValue int64) bool {
	integer, ok := integerLiteral.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("IntegerLiteral type is incorrect. Expected: *ast.IntegerLiteral. Got: %T", integer)
		return false
	}

	if integer.Value != expectedValue {
		t.Fatalf("IntegerLiteral.Value is incorrect. Expected: %d. Got: %d", expectedValue, integer.Value)
		return false
	}

	expectedToken := fmt.Sprintf("%d", expectedValue)
	if integer.TokenLiteral() != expectedToken {
		t.Fatalf("IntegerLiteral.TokenLiteral is incorrect. Expected %q. Got %q", expectedToken, integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, identifierLiteral ast.Expression, expectedValue string) bool {
	identifier, ok := identifierLiteral.(*ast.Identifier)
	if !ok {
		t.Fatalf("IdentifierLiteral type is incorrect. Expected *ast.Identifier. Got: %T", identifierLiteral)
		return false
	}

	if identifier.Value != expectedValue {
		t.Fatalf("IdentifierLiteral.Value is incorrect. Expected: %q. Got %q", expectedValue, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != expectedValue {
		t.Fatalf("IdentifierLiteral.TokenLiteral is incorrect. Expceted: %q. Got: %q", expectedValue, identifier.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, booleanLiteral ast.Expression, expectedValue bool) bool {
	boolean, ok := booleanLiteral.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("BooleanLiteral type is incorrect. Expected: *ast.BoolenLiteral. Got: %T", boolean)
		return false
	}

	if boolean.Value != expectedValue {
		t.Fatalf("BooleanLiteral.Value is incorrect. Expected: %t. Got: %t", expectedValue, boolean.Value)
		return false
	}

	// Book does this slightly differently, however it should still work
	boolToken, _ := strconv.ParseBool(boolean.TokenLiteral())
	if boolToken != expectedValue {
		t.Fatalf("BooleanLiteral.TokenLiteral() is incorrect. Expected: %t. Got: %t", expectedValue, boolToken)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	infix, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("InfixExpression type is incorrect. Expected *ast.InfixExpression. Got: %q", expression)
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}

	if infix.Operator != operator {
		t.Fatalf("InfixExpression.Operator is incorrect. Expected: %q. Got: %q", operator, infix.Operator)
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
}
