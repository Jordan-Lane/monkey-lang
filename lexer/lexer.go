package lexer

import "monkeylang/token"

// TODO: Instead of a string we should use store an io.Reader and the filename. That way we can read a file and
// not load the entire script into memory. This would also allow for better error handling (right now we have pretty much nothing)

// Lexer object
type Lexer struct {
	input        string
	position     int  // points to ch byte in the input string
	readPosition int  // points to the next character in the input string
	char         byte // current character
}

// New - Creates new lexer pointer
func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.inititalizePointers()
	return lexer
}

// NextToken - Returns the next token of the input string
// This is where an enum for the token would be beneificial
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.char {
	case '!':
		if lexer.peekChar() == '=' {
			firstChar := lexer.char
			lexer.readChar()
			literal := string(firstChar) + string(lexer.char)
			tok = token.Token{Type: token.BANG_EQUAL, Literal: literal}
		} else {
			tok = token.NewToken(token.BANG, lexer.char)
		}
	case '=':
		if lexer.peekChar() == '=' {
			firstChar := lexer.char
			lexer.readChar()
			literal := string(firstChar) + string(lexer.char)
			tok = token.Token{Type: token.EQUAL, Literal: literal}
		} else {
			tok = token.NewToken(token.ASSIGN, lexer.char)
		}
	case '+':
		tok = token.NewToken(token.PLUS, lexer.char)
	case '-':
		tok = token.NewToken(token.MINUS, lexer.char)
	case '*':
		tok = token.NewToken(token.STAR, lexer.char)
	case '/':
		tok = token.NewToken(token.SLASH, lexer.char)
	case '<':
		tok = token.NewToken(token.LESS, lexer.char)
	case '>':
		tok = token.NewToken(token.GREATER, lexer.char)
	case ',':
		tok = token.NewToken(token.COMMA, lexer.char)
	case ';':
		tok = token.NewToken(token.SEMICOLON, lexer.char)
	case '(':
		tok = token.NewToken(token.LPAREN, lexer.char)
	case ')':
		tok = token.NewToken(token.RPAREN, lexer.char)
	case '{':
		tok = token.NewToken(token.LBRACE, lexer.char)
	case '}':
		tok = token.NewToken(token.RBRACE, lexer.char)
	case '"':
		tok.Type = token.STRING
		tok.Literal = lexer.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(lexer.char) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookUpIdentifier(tok.Literal)
			return tok
		}
		if isDigit(lexer.char) {
			tok.Literal = lexer.readNumber()
			tok.Type = token.INT
			return tok
		}
		tok = token.NewToken(token.ILLEGAL, lexer.char)
	}
	lexer.readChar()
	return tok
}

func (lexer *Lexer) inititalizePointers() {
	lexer.readChar()
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition++
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.readPosition]
}

func (lexer *Lexer) readIdentifier() string {
	startPosition := lexer.position
	for isLetter(lexer.char) {
		lexer.readChar()
	}
	return lexer.input[startPosition:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	startPosition := lexer.position
	for isDigit(lexer.char) {
		lexer.readChar()
	}
	return lexer.input[startPosition:lexer.position]
}

func (lexer *Lexer) readString() string {
	startPosition := lexer.position + 1
	for {
		lexer.readChar()
		if lexer.char == '"' || lexer.char == 0 {
			break
		}
	}

	return lexer.input[startPosition:lexer.position]
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\n' || lexer.char == '\r' {
		lexer.readChar()
	}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' ||
		char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
