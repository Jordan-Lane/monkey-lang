package parser

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/lexer"
	"monkeylang/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // == or !=
	LESSGREATER // < or >
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -x or !x
	CALL        // func(x + y)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var infixPrecedences = map[token.TokenType]int{
	token.EQUAL:      EQUALS,
	token.BANG_EQUAL: EQUALS,
	token.LESS:       LESSGREATER,
	token.GREATER:    LESSGREATER,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.STAR:       PRODUCT,
	token.SLASH:      PRODUCT,
	token.FUNCTION:   CALL,
	token.LPAREN:     CALL,
}

// Parser ...
type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	// Maps that associated a token to a parsing function
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New - Creates new parser pointer
func New(l *lexer.Lexer) *Parser {
	parser := &Parser{lexer: l, errors: []string{}}

	// Read two tokens so that currToken and nextToken are set
	parser.nextToken()
	parser.nextToken()

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.FALSE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.STAR, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.LESS, parser.parseInfixExpression)
	parser.registerInfix(token.GREATER, parser.parseInfixExpression)
	parser.registerInfix(token.EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.BANG_EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseFunctionCallExpression)

	return parser
}

func (parser *Parser) nextToken() {
	parser.currToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
	return
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
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
	default:
		return parser.parseExpressionStatement()
	}
}

// Parsing General Statements

func (parser *Parser) parseExpressionStatement() ast.Statement {
	statement := &ast.ExpressionStatement{Token: parser.currToken}

	statement.Expression = parser.parseExpression(LOWEST)

	if parser.isPeekTokenType(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.currToken}
	block.Statements = []ast.Statement{}

	parser.nextToken()

	for !parser.isCurrTokenType(token.RBRACE) && !parser.isCurrTokenType(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		parser.nextToken()
	}

	return block
}

// Main Parsing Function. Parsing an Expression based on the inputted precedence

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.currToken.Type]
	if prefix == nil {
		parser.noPrefixParseFnError(parser.currToken.Type)
		return nil
	}
	leftExpression := prefix()

	for !parser.isPeekTokenType(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		parser.nextToken()
		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

// Parsing Prefix vs Infix Expressions

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.currToken,
		Operator: parser.currToken.Literal,
	}

	parser.nextToken()
	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.currToken,
		Operator: parser.currToken.Literal,
		Left:     leftExpression,
	}
	precedence := parser.currPrecedence()

	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

// Parsing Grouped Expression

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	expression := parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		parser.peekError(token.RPAREN)
		return nil
	}

	return expression
}

// Parsing Keywords

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

	parser.nextToken()

	statement.Value = parser.parseExpression(LOWEST)

	if parser.isPeekTokenType(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.currToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	parser.nextToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.isPeekTokenType(token.ELSE) {
		parser.nextToken()

		if !parser.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = parser.parseBlockStatement()
	}

	return expression
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{Token: parser.currToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	functionLiteral.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	functionLiteral.Body = parser.parseBlockStatement()

	return functionLiteral
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}

	if parser.isPeekTokenType(token.RPAREN) {
		parser.nextToken()
		return parameters
	}

	parser.nextToken()

	firstParam := &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
	parameters = append(parameters, firstParam)

	for parser.isPeekTokenType(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		param := &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
		parameters = append(parameters, param)
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (parser *Parser) parseFunctionCallExpression(functionName ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: parser.currToken, Function: functionName}
	expression.Arguments = parser.parseFunctionArguments()
	return expression
}

func (parser *Parser) parseFunctionArguments() []ast.Expression {
	args := []ast.Expression{}

	if parser.isPeekTokenType(token.RPAREN) {
		parser.nextToken()
		return args
	}

	parser.nextToken()
	args = append(args, parser.parseExpression(LOWEST))

	for parser.isPeekTokenType(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		args = append(args, parser.parseExpression(LOWEST))
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (parser *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: parser.currToken}

	parser.nextToken()

	statement.ReturnValue = parser.parseExpression(LOWEST)

	if parser.isPeekTokenType(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

// Parsing Literals

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.currToken}

	value, err := strconv.ParseInt(parser.currToken.Literal, 0, 64)
	if err != nil {
		errorMsg := fmt.Sprintf("Could not parse %q as an integer", literal.Value)
		parser.errors = append(parser.errors, errorMsg)
		return nil
	}

	literal.Value = value
	return literal
}

func (parser *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: parser.currToken, Value: parser.isCurrTokenType(token.TRUE)}
}

// Helper functions

func (parser *Parser) isCurrTokenType(expectedTokenType token.TokenType) bool {
	return parser.currToken.Type == expectedTokenType
}

func (parser *Parser) isPeekTokenType(expectedTokenType token.TokenType) bool {
	return parser.peekToken.Type == expectedTokenType
}

func (parser *Parser) expectPeek(expectedTokenType token.TokenType) bool {
	if parser.isPeekTokenType(expectedTokenType) {
		parser.nextToken()
		return true
	}
	return false
}

func (parser *Parser) peekPrecedence() int {
	if precedence, ok := infixPrecedences[parser.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (parser *Parser) currPrecedence() int {
	if precedence, ok := infixPrecedences[parser.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

// Parsing Errors

func (parser *Parser) Errors() []string {
	return parser.errors
}
func (parser *Parser) peekError(expectedTokenType token.TokenType) {
	errorMsg := fmt.Sprintf("Expected token type %s, got %s instead", expectedTokenType, parser.peekToken.Type)
	parser.errors = append(parser.errors, errorMsg)
	return
}

func (parser *Parser) noPrefixParseFnError(tokenType token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function found for tokentype: %s", tokenType)
	parser.errors = append(parser.errors, msg)
	return
}
