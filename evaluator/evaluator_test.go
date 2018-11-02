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
		if !testIntegerObject(t, evaluated, test.expectedValue) {
			t.Fatalf("Full input: " + test.input)
		}
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
		if !testBooleanObject(t, evaluated, test.expectedBool) {
			t.Fatalf("Full input: " + test.input)
		}
	}

}

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
		if !testBooleanObject(t, evaluated, test.expectedBool) {
			t.Fatalf("Full input: " + test.input)
		}
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
