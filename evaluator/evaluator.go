package evaluator

import (
	"interpreter/ast"
	"interpreter/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        case *ast.Program:
            return evalStatements(node.Statements)
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        case *ast.PrefixExpression:
            right := Eval(node.Right)
            return evalPrefixExpression(node.Operator, right)
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.Boolean:
            return nativeBoolToBooleanObject(node.Value)
        case *ast.InfixExpression:
            left := Eval(node.Left)
            right := Eval(node.Right)
            return evalInfixExpression(node.Operator, left, right)
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

func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
        case "!":
            return evalBangOperatorExpression(right)
        case "-":
            return evalMinuxPrefixOperatorExpression(right)
        default:
            return NULL
    }
}

func evalBangOperatorExpression(right object.Object) object.Object {
    switch right {
        case TRUE:
            return FALSE
        case FALSE:
            return TRUE
        case NULL:
            return TRUE
        default:
            return FALSE
    }
}

func evalMinuxPrefixOperatorExpression(right object.Object) object.Object {
   if right.Type() != object.INTEGER_OBJ {
      return NULL
   }

   value := right.(*object.Integer).Value

   return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
    switch {
        case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
            return evalIntegerInfixExpression(operator, left, right)
        case operator == "==":
            return nativeBoolToBooleanObject(left == right)
        case operator == "!=":
            return nativeBoolToBooleanObject(left != right)
        default:
            return NULL
    }
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.Integer).Value
    rightVal := right.(*object.Integer).Value
    switch operator {
        case "+":
            return &object.Integer{Value: leftVal + rightVal}
        case "-":
            return &object.Integer{Value: leftVal - rightVal}
        case "*":
            return &object.Integer{Value: leftVal * rightVal}
        case "/":
            return &object.Integer{Value: leftVal / rightVal}
        case "<":
            return nativeBoolToBooleanObject(leftVal < rightVal)
        case ">":
            return nativeBoolToBooleanObject(leftVal > rightVal)
        case "==":
            return nativeBoolToBooleanObject(leftVal == rightVal)
        case "!=":
            return nativeBoolToBooleanObject(leftVal != rightVal)
        default:
            return NULL
    }
}
