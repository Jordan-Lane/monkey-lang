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

// BlockStatementStruct - implements Statement Interface
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (blockStatement *BlockStatement) statementNode()       {}
func (blockStatement *BlockStatement) TokenLiteral() string { return blockStatement.Token.Literal }
func (blockStatement *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range blockStatement.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

// ExpressionStatement struct - implements Statement Interface
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expressionStatement *ExpressionStatement) statementNode() {}
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}
func (expressionStatement *ExpressionStatement) String() string {

	// TODO: Remove nil check once expressions are implemented in parser
	if expressionStatement.Expression != nil {
		return expressionStatement.Expression.String()
	}
	return ""
}

// PrefixExpression struct - implements Expression interface
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (prefixExpression *PrefixExpression) expressionNode()      {}
func (prefixExpression *PrefixExpression) TokenLiteral() string { return prefixExpression.Token.Literal }
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

func (infixExpression *InfixExpression) expressionNode()      {}
func (infixExpression *InfixExpression) TokenLiteral() string { return infixExpression.Token.Literal }
func (infixExpression *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(infixExpression.Left.String())
	out.WriteString(" " + infixExpression.Operator + " ")
	out.WriteString(infixExpression.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpression *IfExpression) expressionNode()      {}
func (ifExpression *IfExpression) TokenLiteral() string { return ifExpression.Token.Literal }
func (ifExpression *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ifExpression.Condition.String())
	out.WriteString(ifExpression.Consequence.String())

	if ifExpression.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ifExpression.Alternative.String())
	}

	return out.String()
}

// LetStatement struct - implements Statement Interface
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
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

// ReturnStatement struct - implements Statement Interface
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) statementNode()       {}
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }
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

// Identifier struct
type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) String() string       { return identifier.Value }

// BooleanLiteral struct - implements Expression interface
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (booleanLiteral *BooleanLiteral) expressionNode()      {}
func (booleanLiteral *BooleanLiteral) TokenLiteral() string { return booleanLiteral.Token.Literal }
func (booleanLiteral *BooleanLiteral) String() string       { return booleanLiteral.Token.Literal }

// IntegerLiteral stuct - implements Expression interface
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (integerLiteral *IntegerLiteral) expressionNode()      {}
func (integerLiteral *IntegerLiteral) TokenLiteral() string { return integerLiteral.Token.Literal }
func (integerLiteral *IntegerLiteral) String() string       { return integerLiteral.Token.Literal }
