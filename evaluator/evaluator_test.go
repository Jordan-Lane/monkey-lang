package evaluator

import (
	"monkeylang/lexer"
	"monkeylang/object"
	"monkeylang/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int64
	}{
		{"5", 5},
		{"23", 23},
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
	}

	for _, test := range tests {
		evaluated := runMonkeyLang(test.input)
		testBooleanObject(t, evaluated, test.expectedBool)
	}

}

func runMonkeyLang(input string) object.Object {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expectedValue int64) bool {
	integerObject, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("Object type is incorrect. Expected: *object.Integer. Got: %T", obj)
		return false
	}

	if integerObject.Value != expectedValue {
		t.Fatalf("IntegerObject value is incorrect. Expected: %q. Got: %q", expectedValue, integerObject.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expectedValue bool) bool {
	booleanObject, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("Object type is incorrect. Expected: *object.Boolean. Got: %T", obj)
		return false
	}

	if booleanObject.Value != expectedValue {
		t.Fatalf("BooleanObject value is incorrect. Expected: %t. Got: %t", expectedValue, booleanObject.Value)
		return false
	}

	return true
}
