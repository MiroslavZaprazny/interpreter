package evaluator

import (
	"interpreter/ast"
	"interpreter/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        case *ast.Program:
            return evalStatements(node.Statements)
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.Boolean:
            return &object.Boolean{Value: node.Value}
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
    if input {
        return TRUE
    }

    return FALSE
}
