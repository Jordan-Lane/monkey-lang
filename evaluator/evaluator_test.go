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
		if !testBooleanObject(t, evaluated, test.expectedBool) {
			t.Fatalf("Full input: " + test.input)
		}
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
		if !testBooleanObject(t, evaluated, test.expectedBool) {
			t.Fatalf("Full input: " + test.input)
		}
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
			if !testIntegerObject(t, evaluated, int64(expectedInteger)) {
				t.Fatalf("Full input: " + test.input)
			}
		} else {
			testNullObject(t, evaluated)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object type is incorrect. Expected: *ast.Null. Got %T", obj)
		return false
	}

	return true
}
