package parser

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/lexer"
	"monkeylang/token"
)

// Parser ...
type Parser struct {
	lexer *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string
}

// New - Creates new parser pointer
func New(l *lexer.Lexer) *Parser {
	parser := &Parser{lexer: l, errors: []string{}}

	// Read two tokens so that currToken and nextToken are set
	parser.nextToken()
	parser.nextToken()

	return parser
}

// Errors - get all parser errors
func (parser *Parser) Errors() []string {
	return parser.errors
}

// ParseProgram - Parse the tokens into a collection of statements
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.currToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	case token.IF:
		return parser.parseIfStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() ast.Statement {
	statement := &ast.LetStatement{Token: parser.currToken}

	if !parser.expectPeek(token.IDENT) {
		parser.peekError(token.IDENT)
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		parser.peekError(token.ASSIGN)
		return nil
	}

	// TODO: Currently skipping over the expression
	for !parser.isCurrTokenType(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: parser.currToken}

	// TODO: Currently skipping over the expression
	for !parser.isCurrTokenType(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseIfStatement() ast.Statement {
	return nil
}

func (parser *Parser) nextToken() {
	parser.currToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
	return
}

func (parser *Parser) expectPeek(expectedTokenType token.TokenType) bool {
	if parser.isPeekTokenType(expectedTokenType) {
		parser.nextToken()
		return true
	}
	return false
}

func (parser *Parser) peekError(expectedTokenType token.TokenType) {
	errorMsg := fmt.Sprintf("Expected token type %s, got %s instead", expectedTokenType, parser.peekToken.Type)
	parser.errors = append(parser.errors, errorMsg)
	return
}

func (parser *Parser) isCurrTokenType(expectedTokenType token.TokenType) bool {
	return parser.currToken.Type == expectedTokenType
}

func (parser *Parser) isPeekTokenType(expectedTokenType token.TokenType) bool {
	return parser.peekToken.Type == expectedTokenType
}
