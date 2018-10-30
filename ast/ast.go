package ast

import "monkeylang/token"

// Node base interface
type Node interface {
	TokenLiteral() string
}

// Statement node interface
type Statement interface {
	Node
	statementNode()
}

// Expression node interface
type Expression interface {
	Node
	expressionNode()
}

// Program Node - The root node of the AST
type Program struct {
	Statements []Statement
}

// TokenLiteral implementation for the Program Node
func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement struct
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

// TokenLiteral implementation for Let Statements
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
func (letStatement *LetStatement) statementNode()       {}

// Identifier struct
type Identifier struct {
	Token token.Token
	Value string
}

// TokenLiteral implementation for the Identifier
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) expressionNode()      {}
