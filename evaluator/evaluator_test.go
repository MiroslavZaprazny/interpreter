package evaluator

import (
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"5", 5},
        {"10", 10},
        {"-5", -5},
        {"5 + 5 + 10 + 10", 30},
        {"3 * 5  * 2 / 2", 15},
        {"(1 + 1) * 3", 6},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func TestEvalBooleanExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"true", true},
        {"false", false},
        {"1 < 2", true},
        {"1 > 2", false},
        {"1 > 1", false},
        {"1 == 1", true},
        {"1 == 2", false},
        {"1 != 1", false},
        {"1 != 2", true},
        {"true == true", true},
        {"true == false", false},
        {"false == false", true},
        {"false == true", false},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}

func TestBangOperator(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!!false", false},
        {"!!5", true},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}

func TestReturnStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"return 10; 9;", 10},
        {"return 5+5; 9;", 10},
        {"return 10 * 8; 150;", 80},
        {`
            if(10 > 1) {
                if (10 > 1) {
                    return 10;
                }
                return 1;
            }
        `,
        10,
        },
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func testLetStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"let a = 5;", 5},
        {"let a = 5 * 5;", 25},
        {"let a = 5 + 5;", 10},
    }

    for _, tt := range tests {
        testIntegerObject(t, testEval(tt.input), tt.expected)
    }
}

func TestErrorHandling(t *testing.T) {
    tests := []struct {
        input string
        expectedError string
    } {
        {"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
        {"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
        {"-true;", "unknown operator: -BOOLEAN"},
        {"true + true;", "unknown operator: BOOLEAN + BOOLEAN"},
        {"5; false + true;", "unknown operator: BOOLEAN + BOOLEAN"},
        {"foobar", "identifier not found: foobar"},
    }

    for _, tt := range tests {
        evaluated :=  testEval(tt.input)
        errorObj, ok := evaluated.(*object.Error)
        if !ok {
            t.Errorf("no error object returned got %T", evaluated)
        }

        if errorObj.Message != tt.expectedError {
            t.Errorf("Expected error=%q got=%q", tt.expectedError, errorObj.Message)
        }
    }
}

func testIfElseExpression(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"if (true) { 10 }", 10},
        {"if (false) { 10 }", nil},
        {"if (1) { 10 }", 10},
        {"if (1 < 2) { 10 }", 10},
        {"if (1 > 2) { 10 }", nil},
        {"if (1 < 2) { 10 } else { 20 }", 10},
        {"if (1 > 2) { 10 } else { 20 }", 20},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        integer, ok := tt.expected.(int)
        if ok {
            testIntegerObject(t, evaluated, int64(integer))
        } else {
            testNullObject(t, evaluated)
        }
    }
}

func testEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    env := object.NewEnviroment()

    return Eval(p.ParseProgram(), env)
}

func testNullObject(t *testing.T, object object.Object) bool {
    if object != NULL {
        t.Errorf("Expected type NULL got=%T", object)
        return false
    }

   return true
}

func testBooleanObject(t *testing.T, input object.Object, expected bool) bool {
    result, ok := input.(*object.Boolean)

    if !ok {
        t.Errorf("Expected boolean object got %T (%+v)", input, input)
        return false;
    }

    if result.Value != expected {
        t.Errorf("Expected value %t, got=%t", expected, result.Value)
        return false;
    }

    return true
}

func testIntegerObject(t *testing.T, input object.Object, expected int64) bool {
    result, ok := input.(*object.Integer)

    if !ok {
        t.Errorf("Expected integer object got %T (%+v)", input, input)
        return false;
    }

    if result.Value != expected {
        t.Errorf("Expected value %d, got=%d", expected, result.Value)
        return false;
    }

    return true
}
