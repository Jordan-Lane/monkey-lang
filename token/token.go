package token

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}

//TODO: Refactor this to use enums instead of strings (see: iota)
// Also more token types coming later in the book 
const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// Literals and Identifiers
	IDENT = "IDENT"
	INT = "INT"

	// Operators
	ASSIGN = "="
	PLUS = "+"

	// Punctuation
	COMMA = ','
	SEMICOLON = ';'

	// Brackets
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	LET = "LET"
	FUNCTION = "FUNCTION"
)