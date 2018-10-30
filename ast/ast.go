package ast

import "monkeylang/token"

// Node ...
type Node interface {
	TokenLiteral() string
}

// Statement ...
type Statement interface {
	Node
	statementNode()
}

// Expression ...
type Expression interface {
	Node
	expressionNode()
}

// Program - The root node of the AST
type Program struct {
	Statements []Statement
}

// TokenLiteral ...
func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement - Struct consisting of the individual elements of a let statement
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

// TokenLiteral - Returns the token literal for a let statement
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
func (letStatement *LetStatement) statementNode()       {}

// Identifier ...
type Identifier struct {
	Token token.Token
	Value string
}

// TokenLiteral - Returns the toke
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) expressionNode()      {}
