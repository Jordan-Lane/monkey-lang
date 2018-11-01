package evaluator

import (
	"monkeylang/ast"
	"monkeylang/object"
)

func Eval(node ast.Node) object.Object {
	switch castedNode := node.(type) {
	case *ast.Program:
		return evalStatements(castedNode.Statements)
	case *ast.ExpressionStatement:
		return Eval(castedNode.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: castedNode.Value}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: castedNode.Value}
	}

	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}
