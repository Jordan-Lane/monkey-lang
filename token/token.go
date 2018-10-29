package token

// TokenType - token string identifier (change this to enum)
type TokenType string

// Token - the token struct
type Token struct {
	Type    TokenType
	Literal string
}

//TODO: Refactor this to use enums instead of strings (see: iota)
const (
	ILLEGAL = "ILLEGAL"
	EOF     = ""

	// Literals and Identifiers
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Punctuation
	COMMA     = ","
	SEMICOLON = ";"

	// Brackets
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// NewToken - Create a new Token from a tokenType and char
func NewToken(tokenType TokenType, char byte) Token {
	return Token{Type: tokenType, Literal: string(char)}
}

// LookUpIdentifier - return the correct identifer for the inputted string
func LookUpIdentifier(literal string) TokenType {
	var identifierTokenType TokenType = IDENT
	if val, ok := keywords[literal]; ok {
		identifierTokenType = val
	}
	return identifierTokenType
}
