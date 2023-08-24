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
        {"-10", -10},
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

func testEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)

    return Eval(p.ParseProgram())
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
        t.Errorf("Expected interger object got %T (%+v)", input, input)
        return false;
    }

    if result.Value != expected {
        t.Errorf("Expected value %d, got=%d", expected, result.Value)
        return false;
    }

    return true
}
