package lexer

// Lexer object
type Lexer struct {
	input        string
	position     int  // points to ch byte in the input string
	readPosition int  // points to the next character in the input string
	ch           byte // current character
}

// New pointer to a Lexer Object
func New(input string) *Lexer {
	return &Lexer{input: input}
}

// Reads the next character of the input string and sets the Lexer's current char
func readCharacter(l *Lexer) {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}
