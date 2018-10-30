package parser

import (
	"monkeylang/ast"
	"monkeylang/lexer"
	"monkeylang/token"
)

// Parser ...
type Parser struct {
	lexer *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

// New - Creates new parser pointer
func New(l *lexer.Lexer) *Parser {
	parser := &Parser{lexer: l}

	// Read two tokens so that currToken and nextToken are set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) nextToken() {
	parser.currToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

// ParseProgram ...
func (parser *Parser) ParseProgram() *ast.Program {
	return nil
}
