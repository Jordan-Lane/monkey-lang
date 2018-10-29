package lexer

import (
	"monkeylang/token"
	"testing"
)

// Learning: testing.T is a type passed to Test functions to manage test state and support formatted test
//		logs. Logs are accumulated during execution and dumped to standard output when done.
//		https://golang.org/pkg/testing/#T
func TestNextToken(t *testing.T) {
	input :=
		`let five = 5;
	 	 let ten = 10;
		 
		 let add = fn(x, y) {
			x + y;
		 };

		 let result = add(five, ten);
	 	`

	// Learning: A slice of structs
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, testToken := range tests {
		currentToken := lexer.NextToken()

		if currentToken.Type != testToken.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected %q, got=%q",
				i, testToken.expectedType, currentToken.Type)
		}

		if currentToken.Literal != testToken.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected %q, got=%q",
				i, testToken.expectedLiteral, currentToken.Literal)
		}
	}

}
