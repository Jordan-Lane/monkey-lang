package ast

import (
	"bytes"
	"monkeylang/token"
)

// Node base interface
type Node interface {
	TokenLiteral() string
	String() string
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

// Program Node - Implements the Node Interface. It is the AST root node
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

func (program *Program) String() string {
	var out bytes.Buffer

	for _, statement := range program.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

// LetStatement struct - implements Statement Interface
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode() {}

// TokenLiteral implementation for the LetStatement struct
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }

// String implementation for the LetStatement struct
func (letStatement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(letStatement.TokenLiteral() + " ")
	out.WriteString(letStatement.Name.String())
	out.WriteString(" = ")

	// TODO: Remove nil check once expressions are implemented in parser
	if letStatement.Value != nil {
		out.WriteString(letStatement.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// Identifier struct
type Identifier struct {
	Token token.Token
	Value string
}

// TokenLiteral implementation for the Identifier
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) String() string       { return identifier.Value }

// ReturnStatement struct - implements Statement Interface
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) statementNode() {}

// TokenLiteral implementation for the ReturnStatement struct
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }

// String implementation for the ReturnStatement struct
func (returnStatement *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(returnStatement.TokenLiteral() + " ")

	// TODO: Remove nil check once expressions are implemented in parser
	if returnStatement.ReturnValue != nil {
		out.WriteString(returnStatement.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement struct - implements Statement Interface
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// TokenLiteral implementation for ExpressionStatement
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}

func (expressionStatement *ExpressionStatement) statementNode() {}

func (expressionStatement *ExpressionStatement) String() string {

	// TODO: Remove nil check once expressions are implemented in parser
	if expressionStatement.Expression != nil {
		return expressionStatement.Expression.String()
	}
	return ""
}

// IntegerLiteral stuct - implements Expression interface
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// TokenLiteral implementation for IntegerLiteral
func (integerLiteral *IntegerLiteral) TokenLiteral() string { return integerLiteral.Token.Literal }

func (integerLiteral *IntegerLiteral) expressionNode() {}

func (integerLiteral *IntegerLiteral) String() string { return integerLiteral.Token.Literal }

// PrefixExpression struct - implements Expression interface
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// TokenLiteral implementation for PrefixExpression
func (prefixExpression *PrefixExpression) TokenLiteral() string { return prefixExpression.Token.Literal }

func (prefixExpression *PrefixExpression) expressionNode() {}

func (prefixExpression *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(prefixExpression.Operator)
	out.WriteString(prefixExpression.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression struct - implements Exprssion Interface
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

// TokenLiteral implementation for InfixExpression
func (infixExpression *InfixExpression) TokenLiteral() string { return infixExpression.Token.Literal }

func (infixExpression *InfixExpression) expressionNode() {}

func (infixExpression *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(infixExpression.Left.String())
	out.WriteString(" " + infixExpression.Operator + " ")
	out.WriteString(infixExpression.Right.String())
	out.WriteString(")")

	return out.String()
}
