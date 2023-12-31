package evaluator

import (
	"fmt"
	"interpreter/ast"
	"interpreter/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
)

func Eval(node ast.Node, env *object.Enviroment) object.Object {
    switch node := node.(type) {
        case *ast.Program:
            return evalProgram(node.Statements, env)
        case *ast.ExpressionStatement:
            return Eval(node.Expression, env)
        case *ast.PrefixExpression:
            right := Eval(node.Right, env)
            if isError(right) {
                return right
            }
            return evalPrefixExpression(node.Operator, right)
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.Boolean:
            return nativeBoolToBooleanObject(node.Value)
        case *ast.InfixExpression:
            left := Eval(node.Left, env)
            if isError(left) {
                return left
            }
            right := Eval(node.Right, env)
            if isError(right) {
                return right
            }
            return evalInfixExpression(node.Operator, left, right)
        case *ast.BlockStatement:
            return evalBlockStatement(node, env)
        case *ast.IfExpression:
            return evalIfExpression(node, env)
        case *ast.ReturnStatement:
            val := Eval(node.ReturnValue, env)
            if isError(val) {
                return val
            }
            return &object.RetrunValue{Value: val}
        case *ast.LetStatemet:
            val := Eval(node.Value, env)
            if isError(val) {
                return val
            }
            env.Set(node.Name.Value, val)
        case *ast.Indentifier:
            return evalIdentifier(node, env)
        case *ast.FunctionLiteral:
            return &object.Function{Parameters: node.Parameters, Env: env, Body: node.Body}
        case *ast.CallExpression:
            function := Eval(node.Function, env)
            if isError(function) {
                return function
            }
            args := evalExpressions(node.Arguments, env)

            if len(args) == 1 && isError(args[0]) {
                return args[0]
            }
            return applyFunction(function, args)
        case *ast.StringLiteral:
            return &object.String{Value: node.Value}
    }

    return nil
}

func evalProgram(statements []ast.Statement, env *object.Enviroment) object.Object {
    var result object.Object

    for _, statement := range statements {
        result = Eval(statement, env)

        switch result := result.(type) {
            case *object.RetrunValue:
                return result.Value
            case *object.Error:
                return result
        }
    }

    return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Enviroment) object.Object {
    condition := Eval(ie.Condition, env)
    if isTruthy(condition) {
        return Eval(ie.Consequence,env)
    } else if ie.Alternative != nil {
        return Eval(ie.Alternative, env)
    } else {
        return NULL
    }
}

func isTruthy(obj object.Object) bool {
    switch obj {
        case NULL:
            return true
        case TRUE:
            return true
        case FALSE:
            return false
        default:
            return true
    }
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
            return evalMinusPrefixOperatorExpression(right)
        default:
            return newError("unknown operator: %s%s", operator, right.Type())
    }
}

func evalExpressions(exps []ast.Expression, env *object.Enviroment) []object.Object {
    var result []object.Object

    for _, e := range exps {
        evaluated := Eval(e, env)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }
        result = append(result, evaluated)
    }

    return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Enviroment) object.Object {
    var result object.Object

    for _, statement := range block.Statements {
        result = Eval(statement, env)

        if result != nil {
            rt := result.Type()
            if rt == object.RETURN_VALUE_OBJ || rt == object.ERRROR_OBJ {
                return result
            }
        }
    }

    return result
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
   if right.Type() != object.INTEGER_OBJ {
       return newError("unknown operator: -%s", right.Type())
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
        case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
            return evalStringInfixExpression(operator, left, right)
        case left.Type() != right.Type():
            return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
        default:
            return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
            return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
    if operator != "+" {
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }
    
    leftVal := left.(*object.String).Value
    rightVal := right.(*object.String).Value

    return &object.String{Value: leftVal + rightVal}
}

func evalIdentifier(node *ast.Indentifier, env *object.Enviroment) object.Object {
    val, ok := env.Get(node.Value)
    if !ok {
        return newError(fmt.Sprintf("identifier not found: %s", node.Value))
    }

    return val
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
    function, ok := fn.(*object.Function)
    if !ok {
        return newError("not a function %s", fn.Type())
    }
    extentedEnv := extentedFunctionEnv(function, args)
    evaluated := Eval(function.Body, extentedEnv)

    return unwrapReturnValue(evaluated)
}

func extentedFunctionEnv(fn *object.Function, args []object.Object) *object.Enviroment {
    env := object.NewEnclosedEnviorment(fn.Env)

    for paramIdx, param := range fn.Parameters {
        env.Set(param.Value, args[paramIdx])
    }

    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    if returnValue, ok := obj.(*object.RetrunValue); ok {
        return returnValue
    }

    return obj
}

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.ERRROR_OBJ
    }

    return false
}
