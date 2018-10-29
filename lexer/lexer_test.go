package lexer

import (
	"monkeylang/token"
	"testing"
)

// Learning: testing.T is a type passed to Test functions to manage test state and support formatted test
//		logs. Logs are accumulated during execution and dumped to standard output when done.
//		https://golang.org/pkg/testing/#T
func TestNextToken(t *testing.T) {
	input := "=+(){},;"

	// Learning: A slice of structs
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		currentToken := l.NextToken()

		if currentToken.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expexted %q, got=%q",
				i, tt.expectedType, currentToken.Type)
		}

		if currentToken.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected %q, got=%q",
				i, tt.expectedLiteral, currentToken.Literal)
		}
	}

}
