package evaluator

import (
	"monkeylang/lexer"
	"monkeylang/object"
	"monkeylang/parser"
	"testing"
)

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input        string
		expectedBool bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, test := range tests {
		evaluated := runMonkeyLang(test.input)
		testBooleanObject(t, evaluated, test.expectedBool)
	}
}
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		evaluated := runMonkeyLang(test.input)
		testIntegerObject(t, evaluated, test.expectedValue)
	}
}

func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input        string
		expectedBool bool
	}{
		{"true", true},
		{"false", false},
		{"2 < 3", true},
		{"10 > 100", false},
		{"5 == 5", true},
		{"5 == 20", false},
		{"30 != 1", true},
		{"30 != 30", false},
		{"true == true", true},
		{"false != false", false},
		{"(2 < 3) == true", true},
		{"(1000 < 10) == true", false},
	}

	for _, test := range tests {
		evaluated := runMonkeyLang(test.input)
		testBooleanObject(t, evaluated, test.expectedBool)
	}

}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, test := range tests {
		evaluated := runMonkeyLang(test.input)
		expectedInteger, ok := test.expectedValue.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(expectedInteger))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			` if (10 > 1) {
				 if (10 > 1) {
				   return 10;
			}
			return 1; }
			`,
			10,
		},
	}
	for _, test := range tests {
		evaluated := runMonkeyLang(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Mismatch types: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Mismatch types: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			` if (10 > 1) {
			if (10 > 1) {
		  		return true + false;
			}	
			return 1; 
		  	}`,
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Unknown identifier: foobar",
		},
	}

	for _, test := range tests {
		evaluated := runMonkeyLang(test.input)
		errorObject, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errorObject.Message != test.expectedMessage {
			t.Errorf("Wrong error message. expected=%q, got=%q",
				test.expectedMessage, errorObject.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, test := range tests {
		testIntegerObject(t, runMonkeyLang(test.input), test.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn ( x ) { x + 2; }"

	evaluated := runMonkeyLang(input)
	function, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Object type is incorrect. Expected: *object.Function. Got: %T", evaluated)
	}

	if len(function.Parameters) != 1 {
		t.Fatalf("Number of parameters is incorrect. Expected: 1. Got: %d", len(function.Parameters))
	}

	expectedParam := "x"
	if function.Parameters[0].String() != expectedParam {
		t.Fatalf("Function.Parameter[0] is incorrect. Expected: %q. Got: %q", expectedParam, function.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if function.Body.String() != expectedBody {
		t.Fatalf("Function.Body is incorrect. Expected: %q. Got: %q", expectedBody, function.Body.String())
	}
}

func runMonkeyLang(input string) object.Object {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	env := object.NewEnvironment()

	return Eval(env, program)
}

func testIntegerObject(t *testing.T, obj object.Object, expectedValue int64) bool {
	integerObject, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object type is incorrect. Expected: *object.Integer. Got: %T", obj)
		return false
	}

	if integerObject.Value != expectedValue {
		t.Errorf("IntegerObject value is incorrect. Expected: %d. Got: %d", expectedValue, integerObject.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expectedValue bool) bool {
	booleanObject, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object type is incorrect. Expected: *object.Boolean. Got: %T", obj)
		return false
	}

	if booleanObject.Value != expectedValue {
		t.Fatalf("BooleanObject value is incorrect. Expected: %t. Got: %t", expectedValue, booleanObject.Value)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object type is incorrect. Expected: *ast.Null. Got %T", obj)
		return false
	}

	return true
}
